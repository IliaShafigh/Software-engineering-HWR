package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"log"
)

func mainChatDisplayConfiguration(cStatusC chan *contivity.ChatroomStatus) *widget.List {
	mainChatDisplay := widget.NewList(
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
	return mainChatDisplay
}

func mainInputEntryConfiguration(a fyne.App, cStatusC chan *contivity.ChatroomStatus, input *mainInputEntry) {
	input.SetPlaceHolder("Write a Message")
	input.OnChanged = func(typed string) {
		if len(typed) >= 50 {
			//input.Disable()
			input.SetText(input.Text[:49])
		}
		if input.Text == "/privateDebug" {
			input.SetText("")
			log.Println("DEBUG PRIVATE CHAT")
			go openRealPrivateWin(a, cStatusC, "", "NONAME")
		}
		if input.Text == "/privateChat" {
			input.SetText("")
			log.Println("Private Chat Please")
			go OpenPrivateWin(a, cStatusC)
		}
	}
}

type mainInputEntry struct {
	widget.Entry
	cStatusC chan *contivity.ChatroomStatus
	errorC   chan contivity.ErrorMessage
}

func (e *mainInputEntry) onEnter() {
	if e.Entry.Text == "" {
		return
	}
	contivity.NGMX(e.Entry.Text, e.cStatusC, e.errorC)
	e.Entry.SetText("")
}

func newMainInputEntry(cStatusC chan *contivity.ChatroomStatus, errorC chan contivity.ErrorMessage) *mainInputEntry {
	entry := &mainInputEntry{}
	entry.ExtendBaseWidget(entry)
	entry.cStatusC = cStatusC
	entry.errorC = errorC
	return entry
}

func (e *mainInputEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	}
}
