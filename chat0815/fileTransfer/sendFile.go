package fileTransfer

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// function to implement into file transfer button
func SendFile(addr net.Addr, myWindow fyne.Window) {
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
			sender(addr, fileSize, fileName, data)
		}, myWindow)
	fileDialog.Resize(fyne.NewSize(750, 500))
	fileDialog.Show()
}

// function handling the sending of the file
func sender(addr net.Addr, fileSize string, fileName string, data []byte) {
	connection, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Println("File-Transfer: conn err :", err, addr.String())
		return
	}
	defer connection.Close()
	request := "NFTX"
	request = fillString(request, 70)
	_, err = connection.Write([]byte(request))

	//send size and name first to create placeholder for file
	fmt.Println("Sending file name and size...")
	//Write first 10 bytes to client telling them the filesize
	_, err = connection.Write([]byte(fileSize))
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

// fillString fills str until it has length of toLength
func fillString(str string, toLength int) string {
	for {
		lengthOfString := len(str)
		if lengthOfString < toLength {
			str = str + ":"
			continue
		}
		break
	}
	return str
}
