package tictacgo

import (
	"chat0815/contivity"
	"fyne.io/fyne/v2/container"
)

func DrawAndShowTTG(chatC chan contivity.ChatStorage, indexOCPT int) {
	chats := <-chatC
	pvStatus := <-chats.Private[indexOCPT].PvStatusC

	board := &Board{}

	SaveBoard(board)

	cont := container.NewGridWithColumns(3)
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			var i = NewBoardIcon(r, c, board)
			SaveBoardIcons(i)
			cont.Add(i)
		}
	}

	WhoStartFirstDialog()

	chats.Private[indexOCPT].TabItem.Content = cont
	chats.AppTabs.Refresh()
	chats.Private[indexOCPT].PvStatusC <- pvStatus
	chatC <- chats
}
