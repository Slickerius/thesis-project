package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	app               fyne.App
	conversationsList binding.StringList
	conversationsMap  map[string]*conversation
	isRunning         bool
	mainWindow        fyne.Window
	debug             *log.Logger
	accountName       binding.String
	accountStatus     binding.String
}

func (gui *GUI) Run(jidChan chan string) {
	// Account input GUI Setup
	loginWindow := gui.app.NewWindow("Enter JID")
	gui.debug.Println("Set loginwindow as master")

	loginWindow.SetCloseIntercept(func() {
		gui.debug.Println("Close window button triggered")
		jidChan <- ""
		loginWindow.Close()
	})

	email := widget.NewEntry()
	email.SetPlaceHolder("alice@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "JID", Widget: email, HintText: "Enter JID you want to chat with"},
		},
		OnCancel: func() {
			gui.debug.Println("Cancelled JID Entry")
			jidChan <- ""
			loginWindow.Close()
		},
		OnSubmit: func() {
			gui.debug.Println("Submitted JID entry")
			jidChan <- email.Text
			gui.accountName.Set(email.Text)
			gui.accountStatus.Set("Offline")
			gui.mainWindow.Show()
			loginWindow.Close()
		},
	}

	loginWindow.SetContent(form)
	loginWindow.Resize(fyne.NewSize(300, 100))
	loginWindow.SetFixedSize(true)
	loginWindow.CenterOnScreen()
	loginWindow.Show()

	gui.isRunning = true
	gui.app.Run()
	gui.isRunning = false
}

func (gui *GUI) Quit() {
	gui.app.Quit()
}

func New(debug *log.Logger) *GUI {
	app := app.New()

	// Main Window GUI Setup
	mainWindow := app.NewWindow("XMPP Client")
	mainWindow.SetMaster()

	conversationsList := binding.NewStringList()
	conversationsMap := map[string]*conversation{}

	chatbox := container.NewStack(container.NewCenter(widget.NewLabel("Your conversations will appear here")))
	setChatBox := func(c *conversation) {
		chatbox.Objects = []fyne.CanvasObject{makeChatBox(c)}
		chatbox.Refresh()
	}
	accountName := binding.NewString()
	accountStatus := binding.NewString()
	sidebar := makeSideBar(conversationsList, conversationsMap, setChatBox, accountName, accountStatus)

	split := container.NewHSplit(sidebar, chatbox)
	split.Offset = 0.2

	mainWindow.SetContent(split)
	mainWindow.Resize(fyne.NewSize(1280, 720))

	gui := &GUI{
		app:               app,
		conversationsList: conversationsList,
		conversationsMap:  conversationsMap,
		isRunning:         false,
		mainWindow:        mainWindow,
		debug:             debug,
		accountName:       accountName,
		accountStatus:     accountStatus,
	}

	return gui
}

func (gui *GUI) ShowPasswordPrompt() string {
	if !gui.isRunning {
		return ""
	}
	w := gui.app.NewWindow("Input Password")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	passChan := make(chan string)
	isFinished := false

	w.SetOnClosed(func() {
		if !isFinished {
			passChan <- ""
			isFinished = true
		}
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Password", Widget: password},
		},
		OnCancel: func() {
			passChan <- ""
			isFinished = true
			w.Close()
		},
		OnSubmit: func() {
			passChan <- password.Text
			isFinished = true
			w.Close()
		},
	}

	w.SetContent(form)
	w.Resize(fyne.NewSize(300, 100))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.Show()

	return <-passChan
}
