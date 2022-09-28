package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"sort"
)

func BuildApp(chatC chan contivity.ChatStorage, errorC chan contivity.ErrorMessage) fyne.App {
	a := app.New()

	go manageChatStorage(chatC)
	go manageLogWindow(errorC, a)

	mainWin := a.NewWindow("chat 0815")
	mainWin.Resize(fyne.NewSize(1200, 600))
	mainWin.SetFixedSize(false)
	mainWin.SetMaster()
	chats := <-chatC
	chatC <- chats
	mainWin.SetOnClosed(func() { contivity.GBXX(chats.GcStatusC) })

	newGroupChatTab(chatC, errorC)
	newGroupChatNavigation(chatC)

	chats = <-chatC
	chats.AppTabs = container.NewAppTabs(chats.GroupChat.TabItem)
	chats.AppTabs.OnChanged = func(tab *container.TabItem) {
		if tab.Text == "Group Chat" {
			//chats := <-chatC
			//chats.Navigation.Remove(chats.Navigation.Objects[0])
			//chats.Navigation.Add(chats.GroupChat.Navigation)
			//chatC <- chats
		} else {
			//chats := <-chatC
			//chats.Navigation.Remove(chats.Navigation.Objects[0])
			//chats.Navigation.Add(chats.Private[chats.AppTabs.CurrentTabIndex()-1].Navigation)
			//chatC <- chats
		}
	}
	chats.Navigation = fyne.NewContainerWithLayout(layout.NewMaxLayout(), chats.GroupChat.Navigation)
	content := container.NewHSplit(chats.Navigation, chats.AppTabs)
	chatC <- chats
	content.SetOffset(0.1)
	mainWin.SetContent(content)
	startUpWin := BuildStartUp(chatC, errorC, a, mainWin)
	startUpWin.Show()
	return a
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

func manageChatStorage(chatC chan contivity.ChatStorage) {
	gcStatusC := make(chan *contivity.GroupChatStatus)
	go manageGcStatusC(gcStatusC)

	groupChat := contivity.GroupChat{
		TabItem:    &container.TabItem{},
		Navigation: &fyne.Container{},
		GcStatusC:  gcStatusC,
		Refresh:    make(chan bool),
	}
	chats := contivity.ChatStorage{
		AppTabs:    &container.AppTabs{},
		Navigation: &fyne.Container{},
		GroupChat:  &groupChat,
		Private:    []*contivity.PrivateChat{},
	}
	for {
		chatC <- chats
		chats = <-chatC
		chats.AppTabs.Refresh() //TODO think about not refreshing/only refreshing what is needed
		chats.GroupChat.Navigation.Refresh()
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

func newGroupChatNavigation(chatC chan contivity.ChatStorage) {
	chats := <-chatC
	chatC <- chats
	list := groupChatNavigationConfiguration(chatC, chats.GcStatusC)
	navigation := fyne.NewContainerWithLayout(layout.NewMaxLayout(), list)
	//Save navigation container in chat storage
	chats = <-chatC
	chats.GroupChat.Navigation = navigation
	chatC <- chats
}

func newGroupChatTab(chatC chan contivity.ChatStorage, errorC chan contivity.ErrorMessage) {
	chats := <-chatC
	chatC <- chats
	chatDisplay := newGroupChatDisplayConfiguration(chats.GcStatusC)
	input := newGroupInputEntry(chats.GcStatusC, errorC)

	lowerBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), input)
	air := layout.NewSpacer()
	chatContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(air, lowerBox, air, air), lowerBox, chatDisplay)
	groupChatTab := container.NewTabItem("Group Chat", chatContainer)
	refresh := make(chan bool)
	go manageDisplayRefresh(refresh, chatDisplay)
	// save the refresh chan and Tabitem in GroupChat
	chats = <-chatC
	chats.GroupChat.Refresh = refresh
	chats.GroupChat.TabItem = groupChatTab
	chatC <- chats
}

func manageDisplayRefresh(refresh chan bool, display *widget.List) {
	for {
		<-refresh
		display.Refresh()
	}
}
