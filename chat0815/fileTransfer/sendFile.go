package fileTransfer

import (
	"fmt"
	"io"
	"net"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
)

func sendFile(conn net.Conn, myWindow fyne.Window) {
	//variables for file information
	var data []byte
	var fileName, fileSize string
	//file dialog for selecting file
	fileDialog := dialog.NewFileOpen(
		func(file fyne.URIReadCloser, _ error) {
			fmt.Println("File \"", file.URI().Name(), "\" selected")
			//read file
			data, _ = io.ReadAll(file)
			//read filesize
			fileSize = strconv.FormatInt(int64(len(data)), 10)
			//read filename
			fileName = file.URI().Name()
		}, myWindow)
	fileDialog.Resize(fyne.NewSize(750, 500))
	fileDialog.Show()
	//send bytes to client telling them the filesize
	_, err := conn.Write([]byte(fileSize))
	if err != nil {
		fmt.Println("Error while sending filesize: ", err)
		panic(err)
	}
	//Write bytes to client containing the filename
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Error while sending filename: ", err)
		panic(err)
	}
	//Initialize a buffer for reading parts of the file in
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error while sending file content: ", err)
		panic(err)
	}
	fmt.Println("File has been sent, closing connection!")
}
