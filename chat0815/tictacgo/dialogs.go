package tictacgo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

//TODO: bool or not bool?
func WhoStartFirstDialog() bool {
	dialog.ShowConfirm(
		"Who start first",
		"Do you want to start first?",
		func(b bool) {
			if PlayerStartingConflict() {
				WhichPlayerStart()
			}
			WhoStartFirst(b)
		},
		fyne.CurrentApp().Driver().AllWindows()[0])
	return true
}

func WhichPlayerStart() {
	if gameStatus.WhoStart == 1 {
		dialog.ShowInformation(
			"Who start first",
			"You can't be first. Enemy wants to start first",
			fyne.CurrentApp().Driver().AllWindows()[0])
	} else {
		dialog.ShowInformation(
			"Who start first",
			"You can't be second. Enemy wants to be second",
			fyne.CurrentApp().Driver().AllWindows()[0])
	}
}
