package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

//TODO REFARCTOR TO OWN .go FILE
func OpenPrivateWin(a fyne.App, cStatusC chan *contivity.ChatroomStatus) {
	listWin := a.NewWindow("List")
	listWin.Resize(fyne.NewSize(600, 600))
	listWin.SetFixedSize(true)

	list := widget.NewList(
		func() int {
			cStatus := <-cStatusC
			cStatusC <- cStatus
			//return with -1 because our own userName is in there as well
			return len(cStatus.UserNames) - 1
		},
		func() fyne.CanvasObject {
			return widget.NewButton("Template", func() {})
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			cStatus := <-cStatusC
			cStatusC <- cStatus
			j := 0
			for key, elem := range cStatus.UserNames {
				if j == i && elem != cStatus.UserName {
					obj.(*widget.Button).SetText(elem)
					obj.(*widget.Button).OnTapped = func() {
						listWin.Hide()
						openRealPrivateWin(a, cStatusC, key)
					}
					return
				} else if elem == cStatus.UserName {
					//Dont add to j if we found our own name, because we dont want to add our name
					continue
				}
				j++
			}
		},
	)
	content := container.NewMax(list)
	listWin.SetContent(content)
	listWin.Show()
}

func openRealPrivateWin(a fyne.App, c chan *contivity.ChatroomStatus, key string) {
	privateWin := a.NewWindow("List")
	privateWin.Resize(fyne.NewSize(600, 600))
	privateWin.SetFixedSize(true)

	privateWin.Show()
}
