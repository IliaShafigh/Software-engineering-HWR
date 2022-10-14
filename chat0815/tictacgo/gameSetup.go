package tictacgo

import "fyne.io/fyne/v2/theme"

//set who start first with consideration of enemy pick
func WhoStartFirst(b bool) {
	if b || gameStatus.WhoStart == 0 {
		boardContainer.whoStart = 1
		boardContainer.myturn = true
		gameStatus.WhoStart = 1
	} else {
		boardContainer.whoStart = 0
		boardContainer.myturn = false
		gameStatus.WhoStart = 0
	}
}

func PlayerStartingConflict() bool {

	if boardContainer.whoStart == -1 {
		return false
	}
	if gameStatus.WhoStart == boardContainer.whoStart {
		return true
	}
	return false
}

func (i *BoardIcon) Reset() {
	i.SetResource(theme.ViewFullScreenIcon())
}

func (b *Board) Reset() {
	for i := range b.pieces {
		b.pieces[i][0] = 0
		b.pieces[i][1] = 0
		b.pieces[i][2] = 0
	}

	b.finished = false
	b.turn = 0
	b.whoStart = -1
	gameStatus.WhoStart = -1
	gameStatus.TurnNumber = 0
}
