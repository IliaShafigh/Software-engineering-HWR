package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"log"
	"net"
)

func BuildStartUp(cStatusC chan *contivity.ChatroomStatus, a fyne.App, mainWin fyne.Window) fyne.Window {

	w := a.NewWindow("connect")
	w.Resize(fyne.NewSize(300, 100))
	w.SetFixedSize(true)

	nameEntry := widget.NewEntry()
	nameEntryConfig(nameEntry)
	ipEntry := widget.NewEntry()
	ipEntryConfig(ipEntry)
	//confirm button
	conf := widget.NewButton("Confirm", func() {
		confButtonClicked(cStatusC, nameEntry.Text, ipEntry.Text, mainWin, w)
	})
	content := container.NewVBox(nameEntry, ipEntry, conf)
	w.SetContent(content)
	return w
}

func confButtonClicked(cStatusC chan *contivity.ChatroomStatus, name, ip string, mainWin fyne.Window, w fyne.Window) {
	connIp := net.ParseIP(ip)
	if connIp == nil && ip != "" {
		log.Println("StartUp: wrong format of ip")
	} else if ip == "" {
		w.Hide()
		mainWin.Show()
	} else {
		check := make(chan bool)
		connAddr := contivity.TcpAddr(connIp)
		go contivity.GetStatusUpdate(&connAddr, cStatusC, check)
		if <-check {
			w.Hide()
			mainWin.Show()
		}
	}
}

func ipEntryConfig(entry *widget.Entry) {
	entry.SetPlaceHolder("Enter ip of one participant (empty if your the first one)")
	entry.OnChanged = func(typed string) {
		if len(typed) >= 20 {
			entry.SetText(entry.Text[:14])
		}
	}
}

func nameEntryConfig(entry *widget.Entry) {
	entry.SetPlaceHolder("Enter your nick name")
	entry.OnChanged = func(typed string) {
		if len(typed) >= 20 {
			entry.SetText(entry.Text[:15])
		}
	}
}
