package errPopUps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowLog(logs ErrorMessage, a fyne.App) {
	line1 := widget.NewLabel(logs.Msg)
	line1.Alignment = fyne.TextAlignCenter
	line2 := widget.NewLabel("no error just a message")
	if logs.Err != nil {
		line2.Text = logs.Err.Error()
	}
	line2.Alignment = fyne.TextAlignCenter
	content := container.NewVBox(line1, line2)
	logWin := a.NewWindow("Log Output")
	logWin.SetContent(content)
	logWin.Show()
}

type ErrorMessage struct {
	Err error
	Msg string
}
