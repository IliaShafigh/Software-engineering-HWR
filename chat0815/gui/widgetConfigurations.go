package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"log"
)

func groupChatDisplayConfiguration(cStatusC chan *contivity.GroupChatStatus) *widget.List {
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

type groupInputEntry struct {
	widget.Entry
	cStatusC chan *contivity.GroupChatStatus
	errorC   chan contivity.ErrorMessage
}

func newGroupInputEntry(cStatusC chan *contivity.GroupChatStatus, errorC chan contivity.ErrorMessage) *groupInputEntry {
	entry := &groupInputEntry{}
	entry.ExtendBaseWidget(entry)
	entry.cStatusC = cStatusC
	entry.errorC = errorC

	entry.SetPlaceHolder("Write a Message")
	entry.OnChanged = func(typed string) {
		if len(typed) >= 50 {
			entry.SetText(entry.Text[:49])
		}
		if entry.Text == "/privateDebug" {
			entry.SetText("")
			log.Println("DEBUG PRIVATE CHAT")
			//go openPrivateTab(a, cStatusC, "", "NONAME")
		}
		if entry.Text == "/privateChat" {
			entry.SetText("")
			log.Println("Private Chat Please")
			//go OpenPrivateWin(a, cStatusC)
		}
	}
	return entry
}

func (e *groupInputEntry) onEnter() {
	if e.Entry.Text == "" {
		return
	}
	contivity.NGMX(e.Entry.Text, e.cStatusC, e.errorC)
	e.Entry.SetText("")
}

func (e *groupInputEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	}
}
