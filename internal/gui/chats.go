package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func generateChatsLabel(c *conversation) []fyne.CanvasObject {
	chatsLabel := []fyne.CanvasObject{}
	messages, _ := c.messageList.Get()
	for _, messageItem := range messages {
		messageObj := messageItem.(*message)
		messageLabel := widget.NewLabel(messageObj.content)
		messageLabel.Wrapping = fyne.TextWrapWord
		if messageObj.sent {
			messageLabel.Alignment = fyne.TextAlignTrailing
		}
		messageLabel.Refresh()
		chatsLabel = append(chatsLabel, messageLabel)
	}
	return chatsLabel
}

func makeChatBox(c *conversation) fyne.CanvasObject {
	chatsLabel := generateChatsLabel(c)
	chatsBase := container.NewVBox(chatsLabel...)
	chats := container.NewVScroll(chatsBase)
	chats.ScrollToBottom()

	c.dataListener = binding.NewDataListener(func() {
		chatsBase.Objects = generateChatsLabel(c)
		chatsBase.Refresh()
		chats.ScrollToBottom()
	})

	c.messageList.AddListener(c.dataListener)

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
				sent:    true,
			}
			c.messageList.Append(message)
			entry.SetText("")
			c.latestMessage.Set(message.content)
			// TODO: Implement sent to xmpp client
		}),
	)
	return container.NewPadded(container.NewBorder(nil, nil, nil, toolbar, entry))
}
