package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

// Define the size of how big the chunks of data will be send each time
// can be between 1 to 65495 bytes
const BUFFERSIZE = 1024

func main() {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//getFile2(conn)
	path := fileLocation()
	fileSaver(conn, path)
}

//ort für die datei auswählen und übergeben
func fileLocation() (filePath string) {
	myApp := app.New()
	//New title and window
	myWindow := myApp.NewWindow("Client")
	// resize window
	myWindow.Resize(fyne.NewSize(400, 400))
	button := widget.NewButton("Save File", func() {
		file_Dialog := dialog.NewFolderOpen(
			func(file fyne.ListableURI, _ error) {
				fileFolder := file.Name()
				filePath = "/" + strings.TrimLeft(file.String(), "file://")
				fmt.Println("Ordner der Datei: ", fileFolder)
				fmt.Println("Pfad der Datei: ", filePath)
			}, myWindow)
		file_Dialog.Show()
	})
	myWindow.SetContent(container.NewVBox(
		button,
	))
	myWindow.ShowAndRun()
	return filePath
}

//datei empfangen und an gewünschten ort erstellen
func fileSaver(conn net.Conn, filePath string) {
	fmt.Println("Verbindung hat geklappt, name und size werden empfangen")
	//Create buffer to read in the name and size of the file
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	//Get the filesize
	conn.Read(bufferFileSize)
	//Strip the ':' from the received size, convert it to a int64
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	//Get the filename
	conn.Read(bufferFileName)
	//Strip the ':' once again but from the received file name now
	fileName := strings.Trim(string(bufferFileName), ":")
	//Create a new file to write in
	newFile, err := os.Create(filePath + "/" + fileName)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	//Create a variable to store in the total amount of data that we received already
	//var receivedBytes int64
	//Start writing in the file
	io.CopyN(newFile, conn, fileSize)
	fmt.Println("Received file completely!")
}
