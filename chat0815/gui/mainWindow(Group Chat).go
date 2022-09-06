package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func BuildApp(cStatus contivity.ChatroomStatus, refresh chan bool) fyne.App {
	a := app.New()

	w := a.NewWindow("Hello World")
	w.Resize(fyne.NewSize(1200, 600))
	w.SetFixedSize(true)

	chatDisplay := widget.NewList(
		func() int {
			return len(*cStatus.ChatContent)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			contents := *cStatus.ChatContent

			obj.(*widget.Label).SetText(contents[len(*cStatus.ChatContent)-1-i])
		},
	)
	//Refresh request
	go refreshChatDisplay(refresh, chatDisplay)
	//Input Chat Console
	input := widget.NewEntry()
	inputEntryConfiguration(input)

	//Send Button
	send := widget.NewButton("Send it!", func() {
		sendButtonClicked(input, cStatus, chatDisplay)
	})

	chat := container.NewMax(chatDisplay)
	navigation := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), input, send)
	content := container.NewHSplit(navigation, chat)
	w.SetContent(content)
	w.Show()
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
		input.Resize(fyne.NewSize(390, 100))
		//TODO CAP the length of the length
		if len(typed) >= 43 {
			//input.Disable()
			input.SetText(input.Text[:42])
		}
	}
}

func sendButtonClicked(input *widget.Entry, cStatus contivity.ChatroomStatus, chatDisplay *widget.List) {
	if input.Text == "" {
		return
	}
	contivity.SendMessageToGroup(cStatus, input.Text)
	//remove later just debuggin
	//*cStatus.ChatContent = append(*cStatus.ChatContent, input.Text)
	chatDisplay.Refresh()

	input.SetText("")
}
