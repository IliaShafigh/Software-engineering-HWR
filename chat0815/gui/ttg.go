package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
)

func drawAndShowTTG(chatC chan contivity.ChatStorage, indexOCPT int) {
	//Draw gameStatus
	chats := <-chatC
	pvStatus := <-chats.Private[indexOCPT].PvStatusC
	if pvStatus.Ttg.Running {
		if pvStatus.Ttg.MyTurn {
			//tabItem := drawMyTurn(pvStatus.Ttg)
		} else {

		}
	} else {

	}
	pvStatus.Ttg.Running = true
	pvStatus.Ttg.MyTurn = false
	SendTtgMove(chatC, indexOCPT)
}

func drawMyTurn(ttg *contivity.TTGGameStatus) *container.TabItem {
	//cells := []*fyne.Container{}
	//for _, c := range ttg.GameField {
	//
	//}
	content := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	return container.NewTabItem("GLEICHER NAME", content)
}

//TODO REFACTOR TO contivity and implement functionality
func SendTtgMove(chatC chan contivity.ChatStorage, indexOCPT int) {

}
