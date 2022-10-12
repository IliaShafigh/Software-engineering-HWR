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

// function to implement into client to save file
func saveFile(connection net.Conn, myWindow fyne.Window) {
	//variable to store the destination of the file
	var filePath string
	//file dialog to pick the destination of the file
	fileDialog := dialog.NewFolderOpen(
		func(file fyne.ListableURI, _ error) {
			//get operating system to determine the path format
			osVersion := runtime.GOOS
			switch osVersion {
			case "windows":
				filePath = strings.TrimLeft(file.String(), "file://")
			case "linux":
				filePath = "/" + strings.TrimLeft(file.String(), "file://")
				//TODO: add MAC OS support
			}
			fmt.Println("Selected path:", filePath)
			//function to save the file
			saver(connection, filePath)
		}, myWindow)
	fileDialog.Resize(fyne.NewSize(600, 600))
	fileDialog.Show()
}

// function handling the saving of the file
func saver(connection net.Conn, filePath string) {
	//Create buffer to read in the name and size of the file
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	fmt.Println("Waiting for file name and size...")
	//Get the filesize
	_, err := connection.Read(bufferFileSize)
	if err != nil {
		fmt.Println("Couldn't read file size: ", err)
		panic(err)
	} else {
		fmt.Println("File size received")
	}
	//Strip the ':' from the received size, convert it to an int64
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	//Get the filename
	_, err = connection.Read(bufferFileName)
	if err != nil {
		fmt.Println("Couldn't read file name: ", err)
		panic(err)
	} else {
		fmt.Println("File name received")
	}
	//Strip the ':' once again from the received file name
	fileName := strings.Trim(string(bufferFileName), ":")
	//Create a placeholder file to write into with the name and size of the file
	newFile, err := os.Create(filePath + "/" + fileName)
	if err != nil {
		fmt.Println("Error while creating empty file as placeholder: ", err)
		panic(err)
	}
	//start writing in the file
	_, err = io.CopyN(newFile, connection, fileSize)
	if err != nil {
		fmt.Println("Error while writing in placeholder file: ", err)
		panic(err)
	} else {
		fmt.Println("File received successfully!")
		fmt.Println("Location: ", filePath+"/"+fileName)
	}
}
