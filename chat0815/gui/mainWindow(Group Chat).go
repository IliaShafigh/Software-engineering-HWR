package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func BuildApp(cStatusC chan *contivity.ChatroomStatus, errorC chan contivity.ErrorMessage, refresh chan bool) fyne.App {
	a := app.New()

	go manageLogWindow(errorC, a)

	mainWin := a.NewWindow("chat 0815")
	mainWin.Resize(fyne.NewSize(1200, 600))
	mainWin.SetFixedSize(true)
	mainWin.SetMaster()
	mainWin.SetOnClosed(func() { contivity.GBXX(cStatusC) })
	startUpWin := BuildStartUp(cStatusC, errorC, a, mainWin)
	chatDisplay := mainChatDisplayConfiguration(cStatusC)
	//Manages Refresh requests
	go manageChatDisplayRefresh(refresh, chatDisplay)

	//Input Chat Console
	input := newMainInputEntry(cStatusC, errorC)
	mainInputEntryConfiguration(a, cStatusC, input)
	//Send Button
	send := widget.NewButton("Send it!", func() {
		sendButtonClicked(input, cStatusC, errorC)
	})

	chat := container.NewMax(chatDisplay)
	navigation := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), input, send)
	content := container.NewHSplit(navigation, chat)
	mainWin.SetContent(content)
	startUpWin.Show()
	return a
}

func manageChatDisplayRefresh(refresh chan bool, chatDisplay *widget.List) {
	for {
		check := <-refresh
		if check {
			chatDisplay.Refresh()
		}
	}
}
func sendButtonClicked(input *mainInputEntry, cStatusC chan *contivity.ChatroomStatus, errorC chan contivity.ErrorMessage) {
	if input.Text == "" {
		return
	}
	contivity.NGMX(input.Text, cStatusC, errorC)
	input.SetText("")
}

func manageLogWindow(errorC chan contivity.ErrorMessage, a fyne.App) {
	var logs contivity.ErrorMessage
	for {
		logs = <-errorC
		go showLog(logs, a)
	}
}
