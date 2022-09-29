package gui

import (
	"chat0815/contivity"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
)

//TODO Implement Tabs to conclude
func openPrivateTab(chatC chan contivity.ChatStorage, addr string) {

	privateChat := []string{"This is private Chat"}

	privateChatDisplay := widget.NewList(
		func() int {
			return len(privateChat)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Templat")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(privateChat[len(privateChat)-1-i])
		},
	)
	//TODO DELETE LATETR
	cStatusC := make(chan *contivity.GroupChatStatus)
	privEntry := newPrivEntry(cStatusC)

	privSendButton := widget.NewButton("Send", func() {
		privateChat = append(privateChat, privEntry.Text)
		privEntry.SetText("")
		privateChatDisplay.Refresh()
	})
	lowerBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), privEntry, privSendButton)
	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(layout.NewSpacer(), lowerBox, layout.NewSpacer(), layout.NewSpacer()), lowerBox, privateChatDisplay)
	_ = content
}

type privateEntry struct {
	widget.Entry
	cStatusC chan *contivity.GroupChatStatus
}

func (e *privateEntry) onEnter() {
	fmt.Println(e.Entry.Text)
	e.Entry.SetText("")
	cStatus := <-e.cStatusC
	log.Println(cStatus.ChatContent)
	e.cStatusC <- cStatus
}

func newPrivEntry(cStatusC chan *contivity.GroupChatStatus) *privateEntry {
	entry := &privateEntry{}
	entry.ExtendBaseWidget(entry)
	entry.cStatusC = cStatusC
	return entry
}

func (e *privateEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	default:
		e.Entry.KeyDown(key)
		fmt.Printf("Key %v pressed\n", key.Name)
	}
}
