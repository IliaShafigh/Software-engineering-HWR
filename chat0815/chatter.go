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

	//FYNE STUFF
	a := gui.BuildApp(chatContent)
	a.Run()
	//__________________________________________________________________
	l, err := net.Listen("tcp", "")
	if err != nil {
		log.Println("Listener died")
		log.Fatal(err)
	}
	defer l.Close()
	go contivity.RunServer(l)

	connected := make(chan bool)
	go contivity.ConnectPLS(l.Addr(), connected)

	<-connected
}
