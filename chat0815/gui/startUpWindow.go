package gui

import (
	"chat0815/contivity"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"net"
)

func BuildStartUp(chatC chan contivity.ChatStorage, errorC chan contivity.ErrorMessage, a fyne.App, mainWin fyne.Window) fyne.Window {

	startUpWin := a.NewWindow("configure start up")
	startUpWin.Resize(fyne.NewSize(500, 100))
	startUpWin.SetFixedSize(true)
	nameEntry := widget.NewEntry()
	nameEntryConfig(nameEntry)
	ipEntry := widget.NewEntry()
	ipEntryConfig(ipEntry)
	//confirm button
	conf := widget.NewButton("Confirm", func() {
		confButtonClicked(chatC, errorC, nameEntry.Text, ipEntry.Text, mainWin, startUpWin)
	})
	content := container.NewVBox(nameEntry, ipEntry, conf)
	startUpWin.SetContent(content)
	return startUpWin
}

//TODO MAKE CLEANER
func confButtonClicked(chatC chan contivity.ChatStorage, errorC chan contivity.ErrorMessage, name, ip string, mainWin, startUpWin fyne.Window) {
	ownAddr := contivity.TcpAddr(contivity.GetOutboundIP())
	//save name in cStatus
	chats := <-chatC
	gcStatus := <-chats.GcStatusC
	if name != "" {
		name = fmt.Sprintf("%-6s", name)
		gcStatus.UserName = name
		gcStatus.UserNames[contivity.AddrWithoutPort(ownAddr)] = name
		chats.GcStatusC <- gcStatus
	} else {
		chats.GcStatusC <- gcStatus
		errorC <- contivity.ErrorMessage{Err: nil, Msg: "Please input your nickname"}
		return
	}
	connIp := net.ParseIP(ip)
	if connIp == nil && ip != "" {
		log.Println("StartUp: wrong format of ip")
	} else if ip == "" {
		contivity.PrintCStatus(*gcStatus)
		startUpWin.Hide()
		mainWin.Show()
	} else {
		check := make(chan bool)
		connAddr := contivity.TcpAddr(connIp)
		go func() {
			err := contivity.UXXX(connAddr, chatC, check, errorC)
			if err != nil {
				errorC <- contivity.ErrorMessage{Err: err, Msg: "Could not connect to " + connAddr.String()}
			}
		}()
		if <-check {
			startUpWin.Hide()
			mainWin.Show()
		}
	}
	chatC <- chats
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
		if len(typed) >= 6 {
			entry.SetText(entry.Text[:6])
		}
	}
}
