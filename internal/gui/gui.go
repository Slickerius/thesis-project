package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type GUI struct {
	app fyne.App
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

	sidebar := makeSideBar()
	chatbox := makeChatBox()
	split := container.NewHSplit(sidebar, chatbox)
	split.Offset = 0.2

	mainWindow.SetContent(split)
	mainWindow.Resize(fyne.NewSize(1280, 720))
	mainWindow.Show()

	gui := &GUI{
		app: app,
	}

	return gui
}
