package tictactoeGrid

type player struct {
	number int
	name   string
	symbol string
}

var playerOne player
var playerTwo player

var grid [3][3]int
var winRow [3]int
var winCol [3]int
var winDiag [3]int
var winOppoDiag [3]int

func Init() {

	setupGrid()
	setupPlayer("missing values")
	go validateWinningGrid()

}

func setupGrid() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			grid[i][j] = 0
		}
	}
}

//Defines who is player one and who is player two
func setupPlayer(value string) {
	//setup player
}

func updateGrid(row int, column int, player player) string {
	if grid[row][column] == 0 {
		grid[row][column] = player.number
		updateWinningGrid(row, column)
		return "grid updated successfully"
	} else {
		return "position is already occupied"
	}
}

//WinningLogic taken from this site https://jayeshkawli.ghost.io/tic-tac-toe/
func updateWinningGrid(row int, column int) {
	winRow[row] += 1
	winCol[column] += 1
	if row == column {
		winDiag[row] += 1
	}
	if row+column == 2 {
		winOppoDiag[row] += 1
	}
}

func resetWinningGrid() {
	for i := 0; i < 3; i++ {
		winRow[i] = 0
		winCol[i] = 0
		winDiag[i] = 0
		winOppoDiag[i] = 0
	}
}

//Determines if user wins
func validateWinningGrid() bool {
	for i := 0; i < 3; i++ {
		if winRow[i] == 3 || winCol[i] == 3 {
			return true
		}
		if winDiag[i] == 1 || winOppoDiag[i] == 1 {
			continue
		} else {
			return false
		}
	}
	return true
}
