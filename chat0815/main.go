package main

import (
	"chat0815/contivity"
	"chat0815/errPopUps"
	"chat0815/gui"
	"log"
	"net"
	"os"
)

func main() {
	chatC := make(chan contivity.ChatStorage)
	errorC := make(chan errPopUps.ErrorMessage)
	//CStatus Management

	//__________________________________________________________________
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("Listener died")
		log.Fatal(err)
	}
	defer l.Close()
	//start server
	go contivity.RunServer(l, chatC, errorC)
	//FYNE STUFF
	a := gui.BuildApp(chatC, errorC)
	a.Run()
	os.Exit(0)
}
