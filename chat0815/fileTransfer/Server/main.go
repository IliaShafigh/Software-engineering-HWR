package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"

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
	//file := fileDialogWindow()
	//Create a TCP listener on port 8888
	server, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Starting Server failed:", err)
		os.Exit(1)
	}
	defer server.Close() //case main ends
	fmt.Println("Server started! Waiting for connections...")
	for {
		conn, err := server.Accept() //acccepts coonection
		if err != nil {
			fmt.Println("Server could not established connection:", err)
			os.Exit(1)
		}
		fmt.Println("Client connected") //send file to specific target(conn)
		//go sendFileToClient(conn, file)
		size, name, data := filePicker()
		go fileSender(conn, size, name, data)
	}
}

//offnet Fenster um Datei Auszuwählen und sammelt infos über Date (size, name, data)
func filePicker() (fileSize string, fileName string, data []byte) {
	myApp := app.New()
	//New title and window
	myWindow := myApp.NewWindow("Server")
	// resize window
	myWindow.Resize(fyne.NewSize(600, 600))
	button := widget.NewButton("Open file", func() {
		file_Dialog := dialog.NewFileOpen(
			func(file fyne.URIReadCloser, _ error) {
				fmt.Println("Datei \"", file.Name(), "\" wurde ausgewählt")
				data, _ = ioutil.ReadAll(file)
				fileName = fillString(file.Name(), 64)
				fileSize = fillString(strconv.FormatInt(int64(len(data)), 10), 10)
			}, myWindow)
		file_Dialog.Show()
	})
	myWindow.SetContent(container.NewVBox(
		button,
	))
	myWindow.ShowAndRun()
	return fileSize, fileName, data
}

//sendet dateiinfos an connection
func fileSender(connection net.Conn, fileSize string, fileName string, data []byte) {
	fmt.Println("Sending filename and filesize!")
	//Write first 10 bytes to client telling them the filesize
	connection.Write([]byte(fileSize))
	//Write 64 bytes to client containing the filename
	connection.Write([]byte(fileName))
	//Initialize a buffer for reading parts of the file in
	connection.Write(data)
	fmt.Println("File has been sent, closing connection!")
	return
}

// Füllt übrigen bytes mit ":" auf
// damit client nicht auf fehlende Bytes wwartet
// falls name bzw. size nicht den gewünschten bytes entspricht
// ":" ist ein illegales file name zeichen somit keine fehler
func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
