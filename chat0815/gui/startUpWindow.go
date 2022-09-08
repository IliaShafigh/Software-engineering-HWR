package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"log"
	"net"
)

func BuildStartUp(cStatusC chan *contivity.ChatroomStatus, errorC chan contivity.ErrorMessage, a fyne.App, mainWin fyne.Window) fyne.Window {

	startUpWin := a.NewWindow("connect")
	startUpWin.SetFixedSize(true)
	logWin := a.NewWindow("Log Output")
	go manageLogWindow(errorC, logWin)

	nameEntry := widget.NewEntry()
	nameEntryConfig(nameEntry)
	ipEntry := widget.NewEntry()
	ipEntryConfig(ipEntry)
	//confirm button
	conf := widget.NewButton("Confirm", func() {
		confButtonClicked(cStatusC, errorC, nameEntry.Text, ipEntry.Text, mainWin, startUpWin)
	})
	content := container.NewVBox(nameEntry, ipEntry, conf)
	startUpWin.SetContent(content)
	return startUpWin
}

func manageLogWindow(errorC chan contivity.ErrorMessage, win fyne.Window) {
	for {
		logs := <-errorC
		win.Hide()
		showLog(logs, win)
	}
}

func confButtonClicked(cStatusC chan *contivity.ChatroomStatus, errorC chan contivity.ErrorMessage, name, ip string, mainWin, w fyne.Window) {
	connIp := net.ParseIP(ip)
	if connIp == nil && ip != "" {
		log.Println("StartUp: wrong format of ip")
	} else if ip == "" {
		w.Hide()
		mainWin.Show()
	} else {
		check := make(chan bool)
		connAddr := contivity.TcpAddr(connIp)
		go func() {
			err := contivity.GetStatusUpdate(&connAddr, cStatusC, check, errorC)
			if err != nil {
				errorC <- contivity.ErrorMessage{Err: err, Msg: "Could not connect to " + connAddr.String()}
			}
		}()

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
