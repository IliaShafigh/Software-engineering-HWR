package tictacgo

var boardIconContainer [3][3]*BoardIcon
var boardContainer *Board
var gameStatus *TicTacGoStatus

func SaveGameStatus(status *TicTacGoStatus) {
	gameStatus = status
}

func SaveBoardIcons(icon *BoardIcon) {
	boardIconContainer[icon.row][icon.column] = icon
}

func ReturnBoardIcons(row int, column int) *BoardIcon {
	return boardIconContainer[row][column]
}

func SaveBoard(board *Board) {
	boardContainer = board
}
