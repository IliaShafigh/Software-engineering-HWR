package main

import (
	"chat0815/contivity"
	"chat0815/gui"
	"log"
	"net"
	"os"
)

func main() {
	outboundAddr := contivity.TcpAddr(contivity.GetOutboundIP())
	_ = outboundAddr
	chatContent := make([]string, 0)

	chatContent = append(chatContent, "Take care of each other and watch your drink")
	chatContent = append(chatContent, "Welcome to chat0815")

	cStatus := contivity.ChatroomStatus{
		ChatContent: chatContent,
		UserAddr:    []net.Addr{},
		BlockedAddr: []net.Addr{},
		UserNames:   make(map[string]string), //map[IPADRESSE]Name
		UserName:    "",
	}
	//Fill own information
	cStatus.UserAddr = append(cStatus.UserAddr, outboundAddr)
	//Create channel communication
	refresh := make(chan bool)
	cStatusC := make(chan *contivity.ChatroomStatus)
	errorC := make(chan contivity.ErrorMessage)
	//CStatus Management
	go manageCStatus(&cStatus, cStatusC)

	//__________________________________________________________________
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("Listener died")
		log.Fatal(err)
	}
	defer l.Close()
	//start server
	go contivity.RunServer(l, cStatusC, refresh, errorC)
	//FYNE STUFF
	a := gui.BuildApp(cStatusC, errorC, refresh)
	a.Run()
	os.Exit(0)
}

//Provides the pointer to the current ChatRoomStatus, always waits for a pointer in return.
//Should run in own goroutine
func manageCStatus(cStatus *contivity.ChatroomStatus, c chan *contivity.ChatroomStatus) {
	for {
		c <- cStatus
		cStatus = <-c
	}
}
