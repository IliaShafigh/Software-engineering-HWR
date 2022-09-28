package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"net"
	"strings"
)

func openPrivateTab(chatC chan contivity.ChatStorage, addr string) {
	managePvChatStatusC(chatC, contivity.TcpAddr(net.ParseIP(strings.Split(addr, ":")[0])))
	chats := <-chatC
	chatC <- chats
	indexOfCurrentPrivateTab := len(chats.Private) - 1 // after PvChatStatus Initialization
	newPrivateChatTab(chatC, indexOfCurrentPrivateTab)
	newPrivateChatNavigation(chatC, indexOfCurrentPrivateTab)
	//TODO show tabitem usw	speichern in chats
	chats = <-chatC
	chats.AppTabs.Append(chats.Private[indexOfCurrentPrivateTab].TabItem)
	chats.AppTabs.SelectTabIndex(indexOfCurrentPrivateTab + 1)
	chats.Navigation.Remove(chats.GroupChat.Navigation)
	chats.Navigation.Refresh()
	chats.Navigation.Add(chats.Private[indexOfCurrentPrivateTab].Navigation)
	chats.Navigation.Refresh()
	chatC <- chats
}

func managePvChatStatusC(chatC chan contivity.ChatStorage, remoteAddr net.Addr) {
	pvStatusC := make(chan *contivity.PrivateChatStatus)

	pvChat := contivity.PrivateChat{
		TabItem:    &container.TabItem{},
		PvStatusC:  pvStatusC,
		Refresh:    make(chan bool),
		Navigation: &fyne.Container{},
	}
	chats := <-chatC
	chats.Private = append(chats.Private, &pvChat)
	chatC <- chats
	pvStatus := contivity.InitializePrivateChatRoomStatus(remoteAddr)
	go func() {
		for {
			pvStatusC <- pvStatus
			pvStatus = <-pvStatusC
			//TODO MAYBE REFRESH STUFF?
		}
	}()
}

// indexOfCurrentPrivateTab
func newPrivateChatTab(chatC chan contivity.ChatStorage, indexOCPT int) {
	chats := <-chatC
	chatC <- chats
	chatDisplay := newPrivateChatDisplayConfiguration(chats.Private[indexOCPT].PvStatusC)
	input := newPrivateInputEntry(chats.Private[indexOCPT].PvStatusC)

	gcStatus := <-chats.GcStatusC
	pvStatus := <-chats.Private[indexOCPT].PvStatusC
	title := gcStatus.UserNames[contivity.AddrWithoutPort(pvStatus.UserAddr)]
	chats.Private[indexOCPT].PvStatusC <- pvStatus
	chats.GcStatusC <- gcStatus

	lowerBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), input)
	air := layout.NewSpacer()
	chatContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(air, lowerBox, air, air), lowerBox, chatDisplay)
	privateChatTab := container.NewTabItem(title, chatContainer)
	refresh := make(chan bool)
	go manageDisplayRefresh(refresh, chatDisplay)
	//save the refresh chan and TabItem in PrivateChat
	chats = <-chatC
	chats.Private[indexOCPT].Refresh = refresh
	chats.Private[indexOCPT].TabItem = privateChatTab
	chatC <- chats
}

func newPrivateChatDisplayConfiguration(pvStatusC chan *contivity.PrivateChatStatus) *widget.List {
	privateChatDisplay := widget.NewList(
		func() int {
			pvStatus := <-pvStatusC
			pvStatusC <- pvStatus
			return len(pvStatus.ChatContent)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Templat")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			pvStatus := <-pvStatusC
			contents := pvStatus.ChatContent
			obj.(*widget.Label).SetText(contents[len(contents)-1-i])
			pvStatusC <- pvStatus
		},
	)
	return privateChatDisplay
}

// indexOfCurrentPrivateTab
func newPrivateChatNavigation(chatC chan contivity.ChatStorage, indexOCPT int) {
	//TODO TicTacGo implementation
	ttgButton := widget.NewButton("TTG", func() {

	})
	//TODO Schiffeversenken implementation
	svButton := widget.NewButton("SV", func() {

	})
	//TODO File-Transfer implementation
	ftButton := widget.NewButton("F-T", func() {

	})
	navigation := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), ftButton, ttgButton, svButton)
	chats := <-chatC
	chats.Private[indexOCPT].Navigation = navigation
	chatC <- chats
}

type privateEntry struct {
	widget.Entry
	pvStatusC chan *contivity.PrivateChatStatus
	errorC    chan contivity.ErrorMessage
}

func (e *privateEntry) onEnter() {
	if e.Entry.Text == "" {
		return
	}
	//TODO NEW Funktion for private messages
	//contivity.NPM(e.Entry.Text, e.pvStatusC)
	e.Entry.SetText("")
}

func newPrivateInputEntry(pvStatusC chan *contivity.PrivateChatStatus) *privateEntry {
	entry := &privateEntry{}
	entry.ExtendBaseWidget(entry)
	entry.pvStatusC = pvStatusC
	return entry
}

func (e *privateEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	}
}
