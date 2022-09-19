package gui

import (
	"chat0815/contivity"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
)

//TODO REFARCTOR TO OWN .go FILE
func OpenPrivateWin(a fyne.App, cStatusC chan *contivity.ChatroomStatus) {
	listWin := a.NewWindow("List")
	listWin.Resize(fyne.NewSize(600, 600))
	listWin.SetFixedSize(true)

	list := widget.NewList(
		func() int {
			cStatus := <-cStatusC
			cStatusC <- cStatus
			//return with -1 because our own userName is in there as well
			return len(cStatus.UserNames) - 1
		},
		func() fyne.CanvasObject {
			return widget.NewButton("Template", func() {})
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			cStatus := <-cStatusC
			cStatusC <- cStatus
			j := 0
			for key, elem := range cStatus.UserNames {
				if j == i && elem != cStatus.UserName {
					obj.(*widget.Button).SetText(elem)
					obj.(*widget.Button).OnTapped = func() {
						listWin.Hide()
						openRealPrivateWin(a, cStatusC, key, elem)
					}
					return
				} else if elem == cStatus.UserName {
					//Dont add to j if we found our own name, because we dont want to add our name
					continue
				}
				j++
			}
		},
	)
	content := container.NewMax(list)
	listWin.SetContent(content)
	listWin.Show()
}

func openRealPrivateWin(a fyne.App, cStatusC chan *contivity.ChatroomStatus, addr string, name string) {
	privateWin := a.NewWindow("Private Chat " + name)
	privateWin.Resize(fyne.NewSize(600, 600))
	privateWin.SetFixedSize(true)
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

	privEntry := newPrivEntry(cStatusC)

	privSendButton := widget.NewButton("Send", func() {
		privateChat = append(privateChat, privEntry.Text)
		privEntry.SetText("")
		privateChatDisplay.Refresh()
	})
	lowerBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), privEntry, privSendButton)
	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(layout.NewSpacer(), lowerBox, layout.NewSpacer(), layout.NewSpacer()), lowerBox, privateChatDisplay)
	privateWin.SetContent(content)
	privateWin.Show()
}

type privateEntry struct {
	widget.Entry
	cStatusC chan *contivity.ChatroomStatus
}

func (e *privateEntry) onEnter() {
	fmt.Println(e.Entry.Text)
	e.Entry.SetText("")
	cStatus := <-e.cStatusC
	log.Println(cStatus.ChatContent)
	e.cStatusC <- cStatus
}

func newPrivEntry(cStatusC chan *contivity.ChatroomStatus) *privateEntry {
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
