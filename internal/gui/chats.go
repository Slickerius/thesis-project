package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeChatBox() fyne.CanvasObject {
	chat1 := widget.NewLabel("Hello this is a message from client")
	chat1.Wrapping = fyne.TextWrapWord
	chat1.Alignment = fyne.TextAlignTrailing
	chat1.Refresh()
	chat2 := widget.NewLabel("Hello this is a message from alice")
	chat2.Wrapping = fyne.TextWrapWord
	chat2.Refresh()
	chat3 := widget.NewLabel("And this is another reply from client to showcase text wrapping ")
	chat3.Wrapping = fyne.TextWrapWord
	chat3.Alignment = fyne.TextAlignTrailing
	chat3.Refresh()
	chats := container.NewVScroll(container.NewVBox(chat1, chat2, chat3))

	toolbar := makeToolbar()
	input := makeInput()
	return container.NewBorder(toolbar, input, nil, nil, chats)
}

func makeToolbar() fyne.CanvasObject {
	addressCard := widget.NewCard("alice@example.com", "", nil)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {}),
		widget.NewToolbarAction(theme.MediaVideoIcon(), func() {}),
	)
	return container.NewBorder(nil, nil, nil, toolbar, addressCard)
}

func makeInput() fyne.CanvasObject {
	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("Enter your message here")
	entry.SetMinRowsVisible(2)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MailSendIcon(), func() {}),
	)
	return container.NewPadded(container.NewBorder(nil, nil, nil, toolbar, entry))
}
