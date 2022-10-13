package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"image/color"
)

func drawAndShowTTG(chatC chan contivity.ChatStorage, indexOCPT int) {
	//Draw gameStatus
	chats := <-chatC
	pvStatus := <-chats.Private[indexOCPT].PvStatusC

	cont := container.New(layout.NewGridLayout(3))
	for i, jj := range pvStatus.Ttg.GameField {
		cont.Add(container.New(layout.NewMaxLayout()))
		drawCell(jj, cont.Objects[i].(*fyne.Container))
	}

	if pvStatus.Ttg.Running {
		if pvStatus.Ttg.MyTurn {
			//tabItem := drawMyTurn(pvStatus.Ttg)
		} else {

		}
	} else {

	}
	cont.Refresh()
	chats.Private[indexOCPT].TabItem.Content = cont
	chats.AppTabs.Refresh()
	pvStatus.Ttg.Running = true
	pvStatus.Ttg.MyTurn = false
	chats.Private[indexOCPT].PvStatusC <- pvStatus
	chatC <- chats
	SendTtgMove(chatC, indexOCPT)
}

func drawCell(xy int, cont *fyne.Container) {
	switch xy {
	case 1:
		//draw X Cell
		drawEmptyCell(cont)
		drawX(cont)
	case 2:
		//draw Y Cell
		drawEmptyCell(cont)
		drawO(cont)
	case 0:
		//draw empty Cell
		drawEmptyCell(cont)
	default:
	}
}

func drawO(cont *fyne.Container) {
	circle := &canvas.Circle{
		Position1: fyne.Position{
			X: 0,
			Y: 0,
		},
		Position2: fyne.Position{
			X: cont.Size().Width,
			Y: cont.Size().Height,
		},
		Hidden:      false,
		FillColor:   nil,
		StrokeColor: color.Black,
		StrokeWidth: 2,
	}
	cont.Add(circle)
}

func drawX(cont *fyne.Container) {
	line1 := &canvas.Line{
		Position1: fyne.Position{
			X: 0,
			Y: 0,
		},
		Position2: fyne.Position{
			X: cont.Size().Width,
			Y: cont.Size().Height,
		},
		Hidden:      false,
		StrokeColor: color.Black,
		StrokeWidth: 2,
	}
	line2 := &canvas.Line{
		Position1: fyne.Position{
			X: cont.Size().Width,
			Y: 0,
		},
		Position2: fyne.Position{
			X: 0,
			Y: cont.Size().Height,
		},
		Hidden:      false,
		StrokeColor: color.Black,
		StrokeWidth: 2,
	}
	cont.Add(line1)
	cont.Add(line2)
}

func drawEmptyCell(cont *fyne.Container) {
	lineLeft := &canvas.Line{
		Position1: fyne.Position{
			X: 0,
			Y: 0,
		},
		Position2: fyne.Position{
			X: 0,
			Y: cont.Size().Height,
		},
		Hidden:      false,
		StrokeColor: color.Black,
		StrokeWidth: 2,
	}
	lineLeft.Move(fyne.NewPos(0, 0))
	lineLeft.Resize(fyne.NewSize(0, cont.Size().Height))
	lineTop := &canvas.Line{
		Position1: fyne.Position{
			X: 0,
			Y: 0,
		},
		Position2: fyne.Position{
			X: cont.Size().Width,
			Y: 0,
		},
		Hidden:      false,
		StrokeColor: color.Black,
		StrokeWidth: 2,
	}
	lineRight := &canvas.Line{
		Position1: fyne.Position{
			X: cont.Size().Width,
			Y: 0,
		},
		Position2: fyne.Position{
			X: cont.Size().Width,
			Y: cont.Size().Height,
		},
		Hidden:      false,
		StrokeColor: color.Black,
		StrokeWidth: 2,
	}
	lineBottom := &canvas.Line{
		Position1: fyne.Position{
			X: 0,
			Y: cont.Size().Height,
		},
		Position2: fyne.Position{
			X: cont.Size().Width,
			Y: cont.Size().Height,
		},
		Hidden:      false,
		StrokeColor: color.Black,
		StrokeWidth: 2,
	}
	cont.Add(lineLeft)
	cont.Add(lineTop)
	cont.Add(lineRight)
	cont.Add(lineBottom)
}

func drawMyTurn(ttg *contivity.TTGGameStatus) *container.TabItem {
	//cells := []*fyne.Container{}
	//for _, c := range ttg.GameField {
	//
	//}
	content := container.New(layout.NewGridLayout(3))
	return container.NewTabItem("GLEICHER NAME", content)
}

//TODO REFACTOR TO contivity and implement functionality
func SendTtgMove(chatC chan contivity.ChatStorage, indexOCPT int) {

}
