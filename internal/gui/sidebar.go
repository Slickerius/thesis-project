package gui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeSideBar() fyne.CanvasObject {
	data := make([]string, 100)
	for i := range data {
		data[i] = "Test Item " + strconv.Itoa(i)
	}

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			address := widget.NewLabel("alice@example.com")
			address.TextStyle = fyne.TextStyle{Bold: true}
			address.Refresh()

			latestMessage := widget.NewRichTextWithText("This is a long test message that should be truncated for the sake of good display")
			latestMessage.Truncation = fyne.TextTruncateEllipsis
			latestMessage.Refresh()

			lastInteraction := widget.NewLabel("Oct 24")

			return container.NewBorder(nil, nil, nil, lastInteraction, container.NewVBox(address, latestMessage))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			return
		},
	)

	accountCard := widget.NewCard("kenshin@slickerius.com", "Online", nil)

	return container.NewBorder(accountCard, nil, nil, nil, list)
}
