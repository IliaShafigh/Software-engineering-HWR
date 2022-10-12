package fileTransfer

import (
	"fmt"
	"io"
	"net"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
)

// function to implement into file transfer button
func sendFile(connection net.Conn, myWindow fyne.Window) {
	//file dialog to pick a file
	fileDialog := dialog.NewFileOpen(
		func(file fyne.URIReadCloser, _ error) {
			//read data of selected file
			data, _ := io.ReadAll(file)
			//get file name and fill it with ":" to 64 bytes
			fileName := fillString(file.URI().Name(), 64)
			//get file size and fill it with ":" to 10 bytes
			fileSize := fillString(strconv.FormatInt(int64(len(data)), 10), 10)
			fmt.Println("File \"", file.URI().Name(), "\" selected")
			//send file to client
			sender(connection, fileSize, fileName, data)
		}, myWindow)
	fileDialog.Resize(fyne.NewSize(750, 500))
	fileDialog.Show()
}

// function handling the sending of the file
func sender(connection net.Conn, fileSize string, fileName string, data []byte) {
	//send size and name first to create placeholder for file
	fmt.Println("Sending file name and size...")
	//Write first 10 bytes to client telling them the filesize
	_, err := connection.Write([]byte(fileSize))
	if err != nil {
		fmt.Println("Error while sending filesize: ", err)
		panic(err)
	} else {
		fmt.Println("File size sent!")
	}
	//Write 64 bytes to client containing the filename
	_, err = connection.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Error while sending filename: ", err)
		panic(err)
	} else {
		fmt.Println("File name sent!")
	}
	//send file data
	_, err = connection.Write(data)
	if err != nil {
		fmt.Println("Error while sending file content: ", err)
		panic(err)
	} else {
		fmt.Println("File has been sent successfully!")
	}
}

// function to fill file name and size with ":" to 64 and 10 bytes
func fillString(returnedString string, toLength int) string {
	for {
		lengthOfString := len(returnedString)
		if lengthOfString < toLength {
			returnedString = returnedString + ":"
			continue
		}
		break
	}
	return returnedString
}
