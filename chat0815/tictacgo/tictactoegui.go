package tictacgo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func DrawAndShowTTG(gamestatus *TicTacGoStatus) *fyne.Container {

	board := &Board{}

	SaveBoard(board)
	SaveGameStatus(gamestatus)

	cont := container.NewGridWithColumns(3)
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			var i = NewBoardIcon(r, c, board)
			SaveBoardIcons(i)
			cont.Add(i)
		}
	}

	return cont
}
