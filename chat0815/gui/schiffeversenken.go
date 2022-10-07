package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"math"
	"time"
)

func drawAndShowSV(chatC chan contivity.ChatStorage, indexOCPT int) {
	chats := <-chatC
	cont := container.NewWithoutLayout()

	circle1 := circle(80, 150, 150)
	circle2 := circle(50, 400, 500)
	ship := widget.NewIcon(theme.RadioButtonIcon())
	ship.Resize(fyne.NewSize(30, 30))
	cont.Add(circle1)
	cont.Add(circle2)
	cont.Add(ship)
	o := orbit(ship, circle1, cont)
	o.Start()
	cont.Refresh()
	chats.Private[indexOCPT].TabItem.Content = cont
	chats.AppTabs.Refresh()

	chatC <- chats
	SendTtgMove(chatC, indexOCPT)
}

func orbit(ship *widget.Icon, circle *canvas.Circle, cont *fyne.Container) *fyne.Animation {
	return circleAnimation(ship, cont, float64(circle.Size().Width/2.0), float64((circle.Position2.X-circle.Position1.X)/2.0), float64((circle.Position2.Y-circle.Position1.Y)/2.0))
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
	log.Println(radius)
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
