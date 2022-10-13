package gui

import (
	"chat0815/contivity"
	. "chat0815/tictacgo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

func drawAndShowTTG(chatC chan contivity.ChatStorage, indexOCPT int) {
	refresh := make(chan bool)
	chats := <-chatC
	pvStatus := <-chats.Private[indexOCPT].PvStatusC

	board := &Board{}

	cont := container.NewGridWithColumns(3)
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			var i = NewBoardIcon(r, c, board)
			SaveBoardIcons(i)
			cont.Add(i)
		}
	}

	dialog.ShowConfirm("Which player begins", "Do you want to start first?",
		func(b bool) {
			board.DetermineWhoStartFirst(b)
		},
		fyne.CurrentApp().Driver().AllWindows()[0])

	cont.Refresh()
	go manageGameRefresh(refresh, cont)

	chats.Private[indexOCPT].TabItem.Content = cont
	chats.AppTabs.Refresh()
	chats.Private[indexOCPT].PvStatusC <- pvStatus
	chatC <- chats
}

func manageGameRefresh(refresh chan bool, display *fyne.Container) {
	for {
		<-refresh
		display.Refresh()
	}
}
