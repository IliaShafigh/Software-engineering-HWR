package tictacgo

func SendMove() string {

	return string(ConvertGameStatusReadyToSend(gameStatus))
}

//receives inGame data as string
func ReceiveMove(gameStatus string) {

	ttgameStatus := ConvertEnemyStatusToYours(gameStatus)

	EnemyTapped(ttgameStatus)

}
