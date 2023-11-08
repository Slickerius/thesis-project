package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func makeSideBar(conversationsList binding.StringList, conversationsMap map[string]*conversation) fyne.CanvasObject {

	list := widget.NewListWithData(
		conversationsList,
		func() fyne.CanvasObject {
			address := widget.NewLabel("alice@example.com")
			address.TextStyle = fyne.TextStyle{Bold: true}
			address.Refresh()

			latestMessage := widget.NewLabel("This is a long test message that should be truncated for the sake of good display")
			latestMessage.Truncation = fyne.TextTruncateEllipsis
			latestMessage.Refresh()

			lastInteraction := widget.NewLabel("Oct 24")

			return container.NewBorder(nil, nil, nil, lastInteraction, container.NewVBox(address, latestMessage))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			email := item.(binding.String)
			emailVal, _ := email.Get()
			borderContainer := obj.(*fyne.Container)
			vbox := borderContainer.Objects[0].(*fyne.Container)

			latestMessage := vbox.Objects[1].(*widget.Label)
			address := vbox.Objects[0].(*widget.Label)

			address.Bind(email)
			latestMessage.Bind(conversationsMap[emailVal].latestMessage)
		},
	)

	accountCard := widget.NewCard("kenshin@slickerius.com", "Online", nil)

	return container.NewBorder(accountCard, nil, nil, nil, list)
}
