package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
	"sort"
)

func BuildApp(chatC chan contivity.ChatStorage, errorC chan contivity.ErrorMessage) fyne.App {
	a := app.New()

	go manageChatC(chatC)
	go manageLogWindow(errorC, a)

	mainWin := a.NewWindow("chat 0815")
	mainWin.Resize(fyne.NewSize(1200, 600))
	mainWin.SetFixedSize(true)
	mainWin.SetMaster()
	//mainWin.SetOnClosed(func() { contivity.GBXX(cStatusC) })

	groupChatTab := newGroupChatTab(chatC, errorC)
	tabsContainer := container.NewAppTabs(groupChatTab)
	chats := <-chatC
	chats.AppTabs = tabsContainer
	chatC <- chats
	//TODO LISTE OF USERS
	navigation := newGroupNavigation(chatC)
	content := container.NewHSplit(navigation, tabsContainer)
	content.SetOffset(0.1)
	mainWin.SetContent(content)

	startUpWin := BuildStartUp(chatC, errorC, a, mainWin)
	startUpWin.Show()
	return a
}

func newGroupNavigation(chatC chan contivity.ChatStorage) *fyne.Container {
	list := widget.NewList(
		func() int {
			chats := <-chatC
			gcStatus := <-chats.GcStatusC
			chats.GcStatusC <- gcStatus
			chatC <- chats
			return len(gcStatus.UserNames)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("Template", func() {})
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			chats := <-chatC
			gcStatus := <-chats.GcStatusC
			users := GetSortedKeyMap(gcStatus.UserNames)
			for j, userAddr := range users {
				if j == i {
					obj.(*widget.Button).SetText(gcStatus.UserNames[userAddr])
					obj.(*widget.Button).OnTapped = func() {
						openPrivateTab(chatC, userAddr)
					}
					obj.(*widget.Button).Refresh()
					chats.GcStatusC <- gcStatus
					chatC <- chats
					return
				}
			}
		},
	)
	return container.NewMax(list)
}

//GetSortedKeyMap iterates over the given map and returns a sorted slice of its keys(IP adresses)
func GetSortedKeyMap(names map[string]string) []string {
	keys := []string{}
	for k, _ := range names {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func manageChatC(chatC chan contivity.ChatStorage) {
	gcStatusC := make(chan *contivity.GroupChatStatus)
	go manageGcStatusC(gcStatusC)

	groupChat := contivity.GroupChat{
		TabItem:   nil,
		GcStatusC: gcStatusC,
		Refresh:   nil,
	}
	chats := contivity.ChatStorage{
		AppTabs:   &container.AppTabs{},
		GroupChat: &groupChat,
		Private:   []*contivity.PrivateChat{},
	}
	for {
		chatC <- chats
		chats = <-chatC
		chats.AppTabs.Refresh()
	}
}

func manageGcStatusC(gcStatusC chan *contivity.GroupChatStatus) {
	gcStatus := contivity.InitializeGroupChatRoomStatus()
	for {
		gcStatusC <- gcStatus
		gcStatus = <-gcStatusC
	}

}

func manageLogWindow(errorC chan contivity.ErrorMessage, a fyne.App) {
	var logs contivity.ErrorMessage
	for {
		logs = <-errorC
		go showLog(logs, a)
	}
}

func newGroupChatTab(chatC chan contivity.ChatStorage, errorC chan contivity.ErrorMessage) *container.TabItem {
	chats := <-chatC
	log.Println("HERE WE ARE")
	gcStatusC := chats.GroupChat.GcStatusC

	chatDisplay := groupChatDisplayConfiguration(gcStatusC)
	input := newGroupInputEntry(gcStatusC, errorC)

	lowerBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), input)
	air := layout.NewSpacer()
	chatContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(air, lowerBox, air, air), lowerBox, chatDisplay)
	groupChatTab := container.NewTabItem("Group Chat", chatContainer)
	chats.GroupChat.TabItem = groupChatTab
	refresh := make(chan bool)
	chats.GroupChat.Refresh = refresh
	go manageDisplayRefresh(refresh, chatDisplay)
	chats.GroupChat.Refresh <- true

	chatC <- chats

	return groupChatTab
}

func manageDisplayRefresh(refresh chan bool, display *widget.List) {
	for {
		<-refresh
		display.Refresh()
	}
}
