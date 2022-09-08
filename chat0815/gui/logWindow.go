package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"image/color"
)

func showLog(logs contivity.ErrorMessage, logWin fyne.Window) {
	line1 := canvas.NewText(logs.Msg, &color.RGBA{0xff, 0xff, 0xff, 0xff})
	line1.TextSize = 12
	line1.Alignment = fyne.TextAlignCenter
	line2 := canvas.NewText("no error just a message", &color.RGBA{0xff, 0xff, 0xff, 0xff})
	if logs.Err != nil {
		line2.Text = logs.Err.Error()
	}
	line2.TextSize = 12
	line2.Alignment = fyne.TextAlignCenter
	content := container.NewVBox(line1, line2)
	logWin.SetContent(content)
	logWin.Show()
}
