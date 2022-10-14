package tictacgo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TicTacGoStatus struct {
	Won        bool
	MyTurn     bool
	Row        int
	Column     int
	TurnNumber int
	Start      int //0 = not running, 1 = running, 2 = end
	WhoStart   int // 1 first 0 second -1 not set yet
}

type Board struct {
	pieces   [3][3]uint8
	turn     uint8
	finished bool
	yourturn bool
	whoStart int // 1 first 0 second -1 not set yet
}

type BoardIcon struct {
	widget.Icon
	board       *Board
	row, column int
}

var boardIconContainer [3][3]*BoardIcon
var boardContainer *Board
var gameStatus *TicTacGoStatus

func SaveBoardIcons(icon *BoardIcon) {
	boardIconContainer[icon.row][icon.column] = icon
}
func ReturnBoardIcons(row int, column int) *BoardIcon {
	return boardIconContainer[row][column]
}

func (b *Board) result() uint8 {
	// Check for a win in the diagonal direction from top left to bottom right.
	if b.pieces[0][0] != 0 && b.pieces[0][0] == b.pieces[1][1] && b.pieces[1][1] == b.pieces[2][2] {
		return b.pieces[0][0]
	}

	// Check for a win in the diagonal direction from bottom left to top right.
	if b.pieces[0][2] != 0 && b.pieces[0][2] == b.pieces[1][1] && b.pieces[1][1] == b.pieces[2][0] {
		return b.pieces[0][2]
	}

	for i := range b.pieces {
		// Check for a win in the horizontal direction.
		if b.pieces[i][0] != 0 && b.pieces[i][0] == b.pieces[i][1] && b.pieces[i][1] == b.pieces[i][2] {
			return b.pieces[i][0]
		}

		// Check for a win in the vertical direction.
		if b.pieces[0][i] != 0 && b.pieces[0][i] == b.pieces[1][i] && b.pieces[1][i] == b.pieces[2][i] {
			return b.pieces[0][i]
		}
	}

	return 0
}

func (b *Board) CheckIfWinningConditionIsMet(row, column int) {
	b.pieces[row][column] = b.turn%2 + 1

	if b.turn > 3 {
		winner := b.result()
		if winner == 0 {
			if b.turn == 8 {
				dialog.ShowInformation("It is a tie!", "Nobody has won. Better luck next time.", fyne.CurrentApp().Driver().AllWindows()[0])
				b.finished = true
				b.Reset()
				WhoStartFirstDialog()
			}
			return
		}

		number := string(winner + 48) // Number 1 is ascii #49 and 2 is ascii #50.
		dialog.ShowInformation("Player "+number+" has won!", "Congratulations to player "+number+" for winning.", fyne.CurrentApp().Driver().AllWindows()[0])
		b.finished = true
		b.Reset()
		WhoStartFirstDialog()
	}
}

//erbt von CanvasObject, wenn nicht vorhanden, dann passiert nix
func (i *BoardIcon) Tapped(*fyne.PointEvent) {
	if i.board.pieces[i.row][i.column] != 0 || i.board.finished {
		return
	}

	if i.board.yourturn == true {
		if i.board.whoStart == 1 {
			i.SetResource(theme.RadioButtonIcon())
		} else {
			i.SetResource(theme.CancelIcon())

		}
		i.board.yourturn = false
	}

	i.board.CheckIfWinningConditionIsMet(i.row, i.column)
	i.board.turn++

	i.UpdateConnectionData(gameStatus)
}

func (i *BoardIcon) UpdateConnectionData(t *TicTacGoStatus) {
	t.TurnNumber = int(i.board.turn)
	t.Row = i.row
	t.Column = i.column
	t.MyTurn = true
}

//updates Board, when enemy is tapping
func EnemyTapped(status *TicTacGoStatus) {

	i := ReturnBoardIcons(status.Row, status.Column)

	if i.board.pieces[i.row][i.column] != 0 || i.board.finished {
		return
	}

	if i.board.yourturn == false {
		if i.board.whoStart == 1 {
			i.SetResource(theme.RadioButtonIcon())
		} else {
			i.SetResource(theme.CancelIcon())

		}
		i.board.yourturn = true
	}

	i.board.CheckIfWinningConditionIsMet(i.row, i.column)
	i.board.turn++
}

//only used when creating BoardItems. when created boardItems are type of canvasObject and are tappable
func NewBoardIcon(row, column int, board *Board) *BoardIcon {
	i := &BoardIcon{board: board, row: row, column: column}
	i.SetResource(theme.ViewFullScreenIcon())
	i.ExtendBaseWidget(i)
	return i
}

func SaveBoard(board *Board) {
	boardContainer = board
}
