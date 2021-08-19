// Copyright 2019 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/text/transform"

	"mellium.im/communique/internal/client/event"
	"mellium.im/communique/internal/escape"
	"mellium.im/communique/internal/storage"
	"mellium.im/communique/internal/ui"
	"mellium.im/xmpp/roster"
)

func writeMessage(pane *ui.UI, configPath string, msg event.ChatMessage, notNew bool) error {
	if msg.Body == "" {
		return nil
	}

	historyAddr := msg.From
	arrow := "←"
	if msg.Sent {
		historyAddr = msg.To
		arrow = "→"
	}

	historyLine := fmt.Sprintf("%s %s %s\n", time.Now().UTC().Format(time.RFC3339), arrow, msg.Body)

	history := pane.History()

	j := historyAddr.Bare()
	if pane.ChatsOpen() {
		if item, ok := pane.Roster().GetSelected(); ok && item.Item.JID.Equal(j) {
			// If the message JID is selected and the window is open, write it to the
			// history window.
			_, err := io.WriteString(history, historyLine)
			return err
		}
	}

	// If it's not selected (or the message window is not open), mark the item as
	// unread in the roster.
	// If the message isn't a new one (we're just loading history), skip all this.
	if !msg.Sent && !notNew {
		ok := pane.Roster().MarkUnread(j.String(), msg.ID)
		if !ok {
			// If the item did not exist, create it then try to mark it as unread
			// again.
			pane.UpdateRoster(ui.RosterItem{
				Item: roster.Item{
					JID: j,
					// TODO: get the preferred nickname.
					Name:         j.Localpart(),
					Subscription: "none",
				},
			})
			pane.Roster().MarkUnread(j.String(), msg.ID)
		}
		pane.Redraw()
	}
	return nil
}

func loadBuffer(ctx context.Context, pane *ui.UI, db *storage.DB, configPath string, ev roster.Item, msgID string, logger *log.Logger) error {
	history := pane.History()
	history.SetText("")

	iter := db.QueryHistory(ctx, ev.JID.String(), "")
	for iter.Next() {
		cur := iter.Message()
		if cur.ID != "" && cur.ID == msgID {
			_, err := io.WriteString(history, "─\n")
			if err != nil {
				return err
			}
		}
		err := writeMessage(pane, configPath, cur, true)
		if err != nil {
			err = fmt.Errorf("error writing history: %w", err)
			history.SetText(err.Error())
			logger.Println(err)
			return nil
		}
	}
	if err := iter.Err(); err != nil {
		history.SetText(err.Error())
		err = fmt.Errorf("error querying history for %s: %w", ev.JID, err)
		logger.Println(err)
	}
	history.ScrollToEnd()
	return nil
}

// unreadMarkReader wraps an io.Reader in a new reader that will insert an
// unread marker at the given offset.
func unreadMarkReader(r io.Reader, color tcell.Color, offset int64) io.Reader {
	t := escape.Transformer()

	return io.MultiReader(
		transform.NewReader(io.LimitReader(r, offset), t),
		// This marker is used by the text view UI to draw the unread marker
		strings.NewReader("─\n"),
		transform.NewReader(r, t),
	)
}
