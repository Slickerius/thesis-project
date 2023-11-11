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
	accountCard       *widget.Card
}

type LoginData struct {
	JID  string
	Pass string
}

func (gui *GUI) Run(jidChan chan *LoginData) {
	// Account input GUI Setup
	loginWindow := gui.app.NewWindow("Enter JID")

	loginWindow.SetCloseIntercept(func() {
		gui.debug.Println("Close window button triggered")
		jidChan <- &LoginData{
			JID:  "",
			Pass: "",
		}
		loginWindow.Close()
	})

	email := widget.NewEntry()
	email.SetPlaceHolder("alice@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "JID", Widget: email},
			{Text: "Password", Widget: password},
		},
		OnCancel: func() {
			gui.debug.Println("Cancelled JID Entry")
			jidChan <- &LoginData{
				JID:  "",
				Pass: "",
			}
			loginWindow.Close()
		},
		OnSubmit: func() {
			gui.debug.Println("Submitted JID entry")
			jidChan <- &LoginData{
				JID:  email.Text,
				Pass: password.Text,
			}
			gui.accountCard.SetTitle(email.Text)
			gui.accountCard.SetSubTitle("Offline")
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
	sidebar := makeSideBar(conversationsList, conversationsMap, setChatBox)
	accountCard := sidebar.(*fyne.Container).Objects[1].(*widget.Card)

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
		accountCard:       accountCard,
	}

	return gui
}
