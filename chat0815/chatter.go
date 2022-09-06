package main

import (
	"chat0815/contivity"
	"chat0815/gui"
	"log"
	"net"
)

func main() {
	chatContent := new([]string)

	*chatContent = append(*chatContent, "Take care of each other and watch your drink")
	*chatContent = append(*chatContent, "Welcome to chat0815")

	cStatus := contivity.ChatroomStatus{
		ChatContent: chatContent,
		UserAddr:    &[]net.TCPAddr{},
		BlockedAddr: &[]net.TCPAddr{},
	}

	//__________________________________________________________________
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("Listener died")
		log.Fatal(err)
	}
	defer l.Close()
	refresh := make(chan bool)

	go contivity.RunServer(l, cStatus, refresh)

	go contivity.GetStatusUpdate(l.Addr(), cStatus)
	//FYNE STUFF
	a := gui.BuildApp(cStatus, refresh)
	a.Run()

}
