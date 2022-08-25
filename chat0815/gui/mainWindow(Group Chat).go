package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func BuildApp(chatContent *[]string) fyne.App {
	a := app.New()

	w := a.NewWindow("Hello World")
	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)

	chatDisplay := widget.NewList(
		func() int {
			return len(*chatContent)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			contents := *chatContent

			obj.(*widget.Label).SetText(contents[len(*chatContent)-1-i])
		},
	)
	//Input Chat Console
	input := widget.NewEntry()
	inputEntryConfiguration(input)

	//Send Button
	send := widget.NewButton("Send it!", func() {
		sendButtonClicked(input, chatContent, chatDisplay)
	})

	chat := container.NewMax(chatDisplay)
	navigation := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), input, send)
	content := container.NewHSplit(navigation, chat)
	w.SetContent(content)
	w.Show()
	return a
}

//funktionen: Placeholder, TODO Cap Max Letters
func inputEntryConfiguration(input *widget.Entry) {
	input.SetPlaceHolder("Write a Message")
	input.OnChanged = func(typed string) {
		input.Resize(fyne.NewSize(390, 100))
		//TODO CAP the length of the length
		if len(typed) == 43 {
			//input.Disable()
			input.SetText(input.Text[:42])
		}
	}
}

func sendButtonClicked(input *widget.Entry, chatContent *[]string, chatDisplay *widget.List) {
	if input.Text == "" {
		return
	}
	contivity.SendMessageToGroup(input.Text)
	//remove later just debuggin
	*chatContent = append(*chatContent, input.Text)
	chatDisplay.Refresh()

	input.SetText("")
}
