package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	app               fyne.App
	conversationsList binding.StringList
	conversationsMap  map[string]*conversation
}

func (gui *GUI) Run() {
	gui.app.Run()
}

func (gui *GUI) Quit() {
	gui.app.Quit()
}

func New() *GUI {
	app := app.New()

	mainWindow := app.NewWindow("XMPP Client")

	conversationsList := binding.NewStringList()
	conversationsMap := map[string]*conversation{}

	chatbox := container.NewMax(widget.NewLabel("Your conversations will appear here"))
	setChatBox := func (c *conversation)  {
		chatbox.Objects = []fyne.CanvasObject{makeChatBox(c)}
		chatbox.Refresh()
	}
	sidebar := makeSideBar(conversationsList, conversationsMap, setChatBox)

	split := container.NewHSplit(sidebar, chatbox)
	split.Offset = 0.2

	mainWindow.SetContent(split)
	mainWindow.Resize(fyne.NewSize(1280, 720))
	mainWindow.Show()

	gui := &GUI{
		app:               app,
		conversationsList: conversationsList,
		conversationsMap:  conversationsMap,
	}

	return gui
}
