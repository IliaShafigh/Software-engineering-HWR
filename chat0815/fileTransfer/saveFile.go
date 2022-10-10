package fileTransfer

import (
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
)

func saveFile(connection net.Conn, myWindow fyne.Window) {
	//os filesystem for filepath format
	osVersion := runtime.GOOS
	var filePath string
	//file dialog to select file destination
	fileDialog := dialog.NewFolderOpen(
		func(file fyne.ListableURI, _ error) {
			//format path for file destination
			switch osVersion {
			case "windows":
				fmt.Println("Windows")
				filePath = strings.TrimLeft(file.String(), "file://")
			case "linux":
				fmt.Println("Linux")
				filePath = "/" + strings.TrimLeft(file.String(), "file://")
				//TODO: MAC filesystem structure
			}
		}, myWindow)
	fileDialog.Resize(fyne.NewSize(750, 500))
	fileDialog.Show()
	//Create buffer to read the name and size of the file
	var bufferFileName []byte
	var bufferFileSize []byte
	//receive filesize
	_, err := connection.Read(bufferFileSize)
	if err != nil {
		fmt.Println("Couldn't read filesize: ", err)
		panic(err)
	}
	//set filesize and convert it to an int64
	fileSize, _ := strconv.ParseInt(string(bufferFileSize), 10, 64)
	//Get the filename
	_, err = connection.Read(bufferFileName)
	if err != nil {
		fmt.Println("Couldn't read filename: ", err)
		panic(err)
	}
	//set filename
	fileName := string(bufferFileName)
	//Create a new file as placeholder to write in
	newFile, err := os.Create(filePath + "/" + fileName)
	if err != nil {
		fmt.Println("Error while creating empty file as placeholder: ", err)
		panic(err)
	}
	//Start writing in the  placeholder file
	_, err = io.CopyN(newFile, connection, fileSize)
	if err != nil {
		fmt.Println("Error while writing in placeholder file: ", err)
		panic(err)
	}
}
