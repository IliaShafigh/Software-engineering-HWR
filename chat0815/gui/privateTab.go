package gui

import (
	"chat0815/contivity"
	"chat0815/tictacgo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"net"
	"strings"
)

func openPrivateTab(chatC chan contivity.ChatStorage, addr string, name string, a fyne.Window) {
	exists, _ := checkTabExists(chatC, addr)
	if exists {
		return
	}
	managePvChatStatusC(chatC, contivity.TcpAddr(net.ParseIP(strings.Split(addr, ":")[0])))
	chats := <-chatC
	chatC <- chats
	indexOfCurrentPrivateTab := len(chats.Private) - 1 // after PvChatStatus Initialization
	newPrivateChatTab(chatC, indexOfCurrentPrivateTab)
	newPrivateChatNavigation(chatC, indexOfCurrentPrivateTab, a)

	chats = <-chatC
	chats.AppTabs.Append(chats.Private[indexOfCurrentPrivateTab].TabItem)
	chats.Navigation.Remove(chats.GroupChat.Navigation)
	chats.Navigation.Refresh()
	chats.Navigation.Add(chats.Private[indexOfCurrentPrivateTab].Navigation)
	chats.Navigation.Refresh()
	chatC <- chats
	chats.AppTabs.SelectIndex(indexOfCurrentPrivateTab + 1)
}

//Check if the tab exists and selects it and returns the index
func checkTabExists(chatC chan contivity.ChatStorage, addr string) (bool, int) {
	chats := <-chatC
	gcStatus := <-chats.GroupChat.GcStatusC
	name, exists := gcStatus.UserNames[addr]
	if exists {
		for i, tab := range chats.AppTabs.Items {
			if tab.Text == name {
				chats.GroupChat.GcStatusC <- gcStatus
				chatC <- chats
				chats.AppTabs.SelectIndex(i)
				return true, i
			}
		}
		chats.GroupChat.GcStatusC <- gcStatus
		chatC <- chats
		return false, -1
	}
	chats.GroupChat.GcStatusC <- gcStatus
	chatC <- chats
	log.Println("Error, should not happen, please analyse source code")
	panic("Error, should not happen, please analyse source code. Asked to open a tab to a user which is not included in our GcStatus.Usernames list")
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
	chatDisplay := newPrivateChatDisplayConfiguration(chats.Private[indexOCPT])
	input := newPrivateInputEntry(chatC, chats.Private[indexOCPT], indexOCPT)

	gcStatus := <-chats.GcStatusC
	pvStatus := <-chats.Private[indexOCPT].PvStatusC
	title := gcStatus.UserNames[contivity.AddrWithoutPort(pvStatus.UserAddr)]
	chats.Private[indexOCPT].PvStatusC <- pvStatus
	chats.GcStatusC <- gcStatus

	chatContainer := chatCont(input, chatDisplay)
	privateChatTab := container.NewTabItem(title, chatContainer)

	//save the  TabItem in PrivateChat
	chats = <-chatC
	chats.Private[indexOCPT].TabItem = privateChatTab
	chatC <- chats
}

func chatCont(input *privateEntry, chatDisplay *widget.List) *fyne.Container {
	lowerBox := container.New(layout.NewVBoxLayout(), input)
	air := layout.NewSpacer()
	chatContainer := container.New(layout.NewBorderLayout(air, lowerBox, air, air), lowerBox, chatDisplay)
	return chatContainer
}

func newPrivateChatDisplayConfiguration(pvChat *contivity.PrivateChat) *widget.List {
	privateChatDisplay := widget.NewList(
		func() int {
			pvStatus := <-pvChat.PvStatusC
			pvChat.PvStatusC <- pvStatus
			return len(pvStatus.ChatContent)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Templat")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			pvStatus := <-pvChat.PvStatusC
			contents := pvStatus.ChatContent
			msg := contents[len(contents)-1-i]
			obj.(*widget.Label).SetText(msg)
			//Color the label according to msg prefix
			pvChat.PvStatusC <- pvStatus
			if msg[:6] != pvChat.Text {
				if i+1 == len(contents) {
					// last entry is system message, so it has center allign
					obj.(*widget.Label).Alignment = fyne.TextAlignCenter
				} else {
					obj.(*widget.Label).Alignment = fyne.TextAlignTrailing
					obj.(*widget.Label).SetText(strings.Split(msg, ":")[1])
				}
			} else {
				obj.(*widget.Label).Alignment = fyne.TextAlignLeading
			}
		},
	)
	refresh := make(chan bool)
	go manageDisplayRefresh(refresh, privateChatDisplay)
	pvChat.Refresh = refresh
	return privateChatDisplay
}

// indexOfCurrentPrivateTab
func newPrivateChatNavigation(chatC chan contivity.ChatStorage, indexOCPT int, a fyne.Window) {
	chatButton := widget.NewButton("CHAT", func() {
		chats := <-chatC
		chats.Private[indexOCPT].TabItem.Content = chatCont(newPrivateInputEntry(chatC, chats.Private[indexOCPT], indexOCPT), newPrivateChatDisplayConfiguration(chats.Private[indexOCPT]))
		chatC <- chats
	})
	//TODO TicTacGo implementation
	ttgButton := widget.NewButton("TTG", func() {
		tictacgo.DrawAndShowTTG(chatC, indexOCPT)
	})
	//TODO Schiffeversenken implementation
	svButton := widget.NewButton("SV", func() {
		drawAndShowSV(chatC, indexOCPT)
	})
	//TODO File-Transfer implementation
	ftButton := widget.NewButton("F-T", func() {
		chats := <-chatC
		pvStatus := <-chats.Private[indexOCPT].PvStatusC
		ipRemote := pvStatus.UserAddr
		chats.Private[indexOCPT].PvStatusC <- pvStatus
		chatC <- chats
		testFunctionFileTransfer(ipRemote, a)
	})
	navigation := container.New(layout.NewVBoxLayout(), chatButton, ftButton, ttgButton, svButton)
	chats := <-chatC
	chats.Private[indexOCPT].Navigation = navigation
	chatC <- chats
}

type privateEntry struct {
	widget.Entry
	pvStatusC chan *contivity.PrivateChatStatus //TODO since we have chatC and indexOCPT
	chatC     chan contivity.ChatStorage
	indexOCPT int
	errorC    chan contivity.ErrorMessage
}

func (e *privateEntry) onEnter() {
	if e.Entry.Text == "" {
		return
	}
	contivity.NPMX(e.Entry.Text, e.pvStatusC, e.errorC)
	chats := <-e.chatC
	gcStatus := <-chats.GcStatusC
	name := gcStatus.UserName
	chats.GcStatusC <- gcStatus
	e.chatC <- chats
	contivity.AddPrivateMessage(e.Entry.Text, e.chatC, e.indexOCPT, name)
	e.Entry.SetText("")
}

func newPrivateInputEntry(chatC chan contivity.ChatStorage, pvChat *contivity.PrivateChat, indexOCPT int) *privateEntry {
	entry := &privateEntry{}
	entry.ExtendBaseWidget(entry)
	entry.pvStatusC = pvChat.PvStatusC
	entry.chatC = chatC
	entry.indexOCPT = indexOCPT
	return entry
}

func (e *privateEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	}
}
