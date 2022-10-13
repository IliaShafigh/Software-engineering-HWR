package tictacgo

import (
	"chat0815/contivity"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

type privateSaveLocation struct {
	pvStatusC  chan *contivity.PrivateChatStatus
	chatC      chan contivity.ChatStorage
	gameStatus *contivity.TTGGameStatus
	indexOCPT  int
	errMessage chan contivity.ErrorMessage
}

type Board struct {
	pieces     [3][3]uint8
	turn       uint8
	finished   bool
	yourturn   bool
	yoursymbol bool // true = X & false = O
}

type BoardIcon struct {
	widget.Icon
	board       *Board
	row, column int
}

var boardIconContainer [3][3]*BoardIcon
var connectionData privateSaveLocation

func SaveConnectionData(chatC chan contivity.ChatStorage, indexOCPT int) {
	connectionData.chatC = chatC
	chats := <-chatC
	status := <-chats.Private[indexOCPT].PvStatusC
	connectionData.pvStatusC = chats.Private[indexOCPT].PvStatusC
	connectionData.gameStatus = status.Ttg
}

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
				dialog.ShowConfirm("Restart", "Do you want to Restart?", nil, fyne.CurrentApp().Driver().AllWindows()[0])
			}
			return
		}

		number := string(winner + 48) // Number 1 is ascii #49 and 2 is ascii #50.
		dialog.ShowInformation("Player "+number+" has won!", "Congratulations to player "+number+" for winning.", fyne.CurrentApp().Driver().AllWindows()[0])
		b.finished = true
		dialog.ShowConfirm("Restart", "Do you want to Restart?", nil, fyne.CurrentApp().Driver().AllWindows()[0])
	}
}

func (b *Board) Reset() {
	for i := range b.pieces {
		b.pieces[i][0] = 0
		b.pieces[i][1] = 0
		b.pieces[i][2] = 0
	}

	b.finished = false
	b.turn = 0
}

//erbt von CanvasObject, wenn nicht vorhanden, dann passiert nix
func (i *BoardIcon) Tapped(*fyne.PointEvent) {
	if i.board.pieces[i.row][i.column] != 0 || i.board.finished {
		return
	}

	if i.board.yourturn == true {
		if i.board.yoursymbol == true {
			i.SetResource(theme.RadioButtonIcon())
		} else {
			i.SetResource(theme.CancelIcon())

		}
		i.board.yourturn = false
	}

	i.board.CheckIfWinningConditionIsMet(i.row, i.column)
	i.board.turn++

	UpdateConnectionData(i)

	UpdateEnemyBoard()
}

func UpdateConnectionData(i *BoardIcon) {
	connectionData.gameStatus.TurnNumber = int(i.board.turn)
	connectionData.gameStatus.Row = i.row
	connectionData.gameStatus.Column = i.column
	connectionData.gameStatus.MyTurn = true
}

//updates Board, when enemy is tapping
func EnemyTapped(status *contivity.TTGGameStatus) {

	i := ReturnBoardIcons(status.Row, status.Column)

	if i.board.pieces[i.row][i.column] != 0 || i.board.finished {
		return
	}

	if i.board.yourturn == false {
		if i.board.yoursymbol == true {
			i.SetResource(theme.RadioButtonIcon())
		} else {
			i.SetResource(theme.CancelIcon())

		}
		i.board.yourturn = true
	}

	i.board.CheckIfWinningConditionIsMet(i.row, i.column)
	i.board.turn++
}

func (i *BoardIcon) Reset() {
	i.SetResource(theme.ViewFullScreenIcon())
}

//only used when creating BoardItems. when created boardItems are type of canvasObject and are tappable
func NewBoardIcon(row, column int, board *Board) *BoardIcon {
	i := &BoardIcon{board: board, row: row, column: column}
	i.SetResource(theme.ViewFullScreenIcon())
	i.ExtendBaseWidget(i)
	return i
}

//Methods only for Board struct
func (b *Board) DetermineWhoStartFirst(playerStartFirst bool) {
	if playerStartFirst {
		b.yoursymbol = true
		b.yourturn = true
	} else {
		b.yoursymbol = false
	}
}

func UpdateEnemyBoard() {

	packedBoardInformation := ConvertGameStatusReadyToSend(connectionData.gameStatus)

	contivity.GSUX(string(packedBoardInformation), connectionData.pvStatusC, connectionData.errMessage)
}

func UpdateBoard(gameStatus string, ttgameStatus *contivity.TTGGameStatus) {

	ttgameStatus = ConvertEnemyStatusToYours(gameStatus)

	EnemyTapped(ttgameStatus)

}

//StructToJsonConversionTTGGameStatus
func ConvertGameStatusReadyToSend(status *contivity.TTGGameStatus) []byte {
	j, err := json.Marshal(status)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	return j
}

//JsonToStructConversionTTGGameStatus
func ConvertEnemyStatusToYours(status string) *contivity.TTGGameStatus {
	var enemyGameStatus *contivity.TTGGameStatus
	err := json.Unmarshal([]byte(status), &enemyGameStatus)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	return enemyGameStatus
}
