package tictacgo

import (
	"encoding/json"
	"log"
)

//StructToJsonConversionTTGGameStatus
func ConvertGameStatusReadyToSend(status *TicTacGoStatus) []byte {
	j, err := json.Marshal(status)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	return j
}

//JsonToStructConversionTTGGameStatus
func ConvertEnemyStatusToYours(status string) *TicTacGoStatus {
	var enemyGameStatus *TicTacGoStatus
	err := json.Unmarshal([]byte(status), &enemyGameStatus)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	return enemyGameStatus
}
