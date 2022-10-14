package tictacgo

func SendMove(status *TTGGameStatus) string {

	return string(ConvertGameStatusReadyToSend(status))
}

//receives inGame data as string
func ReceiveMove(gameStatus string) {

	ttgameStatus := ConvertEnemyStatusToYours(gameStatus)

	EnemyTapped(ttgameStatus)

}
