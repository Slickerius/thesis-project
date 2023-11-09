package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeChatBox(c *conversation) fyne.CanvasObject {
	chats := widget.NewListWithData(
		c.messageList,
		func() fyne.CanvasObject {
			chat := widget.NewLabel("This is a sample chat message")
			chat.Wrapping = fyne.TextWrapWord
			chat.Refresh()
			return chat
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			messageItem, _ := item.(binding.Untyped).Get()
			message := messageItem.(*message)
			chat := obj.(*widget.Label)
			chat.SetText(message.content)
			if message.sent {
				chat.Alignment = fyne.TextAlignTrailing
				chat.Refresh()
			}
		},
	)

	toolbar := makeToolbar(c)
	input := makeInput(c)
	return container.NewBorder(toolbar, input, nil, nil, chats)
}

func makeToolbar(c *conversation) fyne.CanvasObject {
	addressCard := widget.NewCard(c.email, "", nil)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {}),
		widget.NewToolbarAction(theme.MediaVideoIcon(), func() {}),
	)
	return container.NewBorder(nil, nil, nil, toolbar, addressCard)
}

func makeInput(c *conversation) fyne.CanvasObject {
	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("Enter your message here")
	entry.SetMinRowsVisible(2)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MailSendIcon(), func() {
			message := &message{
				content: entry.Text,
				sent: true,
			}
			c.messageList.Append(message)
			// TODO: Implement sent to xmpp client
		}),
	)
	return container.NewPadded(container.NewBorder(nil, nil, nil, toolbar, entry))
}
