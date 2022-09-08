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
	w := a.NewWindow("chat 0815")
	w.Resize(fyne.NewSize(1200, 600))
	w.SetFixedSize(true)
	w.SetMaster()
	w.SetOnClosed(func() { contivity.SayGoodBye(cStatusC) })
	startUpWin := BuildStartUp(cStatusC, errorC, a, w)
	chatDisplay := widget.NewList(
		func() int {
			cStatus := <-cStatusC
			cStatusC <- cStatus
			return len(cStatus.ChatContent)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			cStatus := <-cStatusC
			contents := cStatus.ChatContent
			obj.(*widget.Label).SetText(contents[len(contents)-1-i])
			cStatusC <- cStatus
		},
	)
	//Refresh request
	go refreshChatDisplay(refresh, chatDisplay)
	//Input Chat Console
	input := widget.NewEntry()
	inputEntryConfiguration(input)

	//Send Button
	send := widget.NewButton("Send it!", func() {
		sendButtonClicked(input, cStatusC, errorC)
	})

	chat := container.NewMax(chatDisplay)
	navigation := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), input, send)
	content := container.NewHSplit(navigation, chat)
	w.SetContent(content)
	startUpWin.Show()
	return a
}

func refreshChatDisplay(refresh chan bool, chatDisplay *widget.List) {
	for {
		check := <-refresh
		if check {
			chatDisplay.Refresh()
		}
	}
}

//funktionen: Placeholder, TODO Cap Max Letters
func inputEntryConfiguration(input *widget.Entry) {
	input.SetPlaceHolder("Write a Message")
	input.OnChanged = func(typed string) {
		if len(typed) >= 43 {
			//input.Disable()
			input.SetText(input.Text[:42])
		}
	}
}

func sendButtonClicked(input *widget.Entry, cStatusC chan *contivity.ChatroomStatus, errorC chan contivity.ErrorMessage) {
	if input.Text == "" {
		return
	}
	contivity.SendMessageToGroup(input.Text, cStatusC, errorC)
	input.SetText("")
}
