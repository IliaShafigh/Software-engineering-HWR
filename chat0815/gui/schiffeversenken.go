package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math"
	"time"
)

func drawAndShowSV(chatC chan contivity.ChatStorage, indexOCPT int) {
	chats := <-chatC
	cont := container.NewWithoutLayout()

	circle1 := circle(80, 150, 150)
	circle2 := &canvas.Circle{
		Position1:   fyne.Position{X: 320, Y: 20},
		Position2:   fyne.Position{X: 440, Y: 140},
		Hidden:      false,
		FillColor:   color.Black,
		StrokeColor: nil,
		StrokeWidth: 0,
	}
	ship := widget.NewIcon(theme.RadioButtonIcon())
	ship.Resize(fyne.NewSize(30, 30))
	cont.Add(circle1)
	cont.Add(circle2)
	cont.Add(ship)
	shipMovement := circleAnimation(ship, cont, 90.0, 150.0, 150.0)
	shipMovement.Start()
	cont.Refresh()
	chats.Private[indexOCPT].TabItem.Content = cont
	chats.AppTabs.Refresh()

	chatC <- chats
	SendTtgMove(chatC, indexOCPT)
}

func circle(radius, centerX, centerY float32) *canvas.Circle {
	return &canvas.Circle{
		Position1:   fyne.Position{X: centerX - radius, Y: centerY - radius},
		Position2:   fyne.Position{X: centerX + radius, Y: centerY + radius}, //Mittelpunkt (160,160), radius = 60 -> radius der bewegung maybe 90
		Hidden:      false,
		FillColor:   color.Black,
		StrokeColor: color.Black,
		StrokeWidth: 0,
	}
}

func circleAnimation(obj *widget.Icon, objHost *fyne.Container, radius, centerX, centerY float64) *fyne.Animation {
	circleA := fyne.NewAnimation(time.Second*5, func(f float32) {
		var x float32
		var y float32

		y = float32(radius*math.Sin(float64(f*2*math.Pi)) + centerY)
		x = float32(radius*math.Cos(float64(f*2*math.Pi)) + centerX)

		obj.Move(fyne.NewPos(x-(obj.Size().Width/2), y-(obj.Size().Height/2))) //Offset so obj center is taken into account
		objHost.Refresh()
	})
	circleA.RepeatCount = -1
	circleA.Curve = fyne.AnimationLinear
	return circleA
}
