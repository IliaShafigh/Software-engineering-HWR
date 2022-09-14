package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// Define the size of how big the chunks of data will be send each time
// can be between 1 to 65495 bytes
const BUFFERSIZE = 1024

func main() {
	conn, err := net.Dial("tcp", "172.20.10.3:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	getFile(conn)
}

func getFile(conn net.Conn) {
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
	newFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	//Create a variable to store in the total amount of data that we received already
	var receivedBytes int64
	//Start writing in the file
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			//Empty the remaining bytes that we don't need from the network buffer
			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			//We are done writing the file, break out of the loop
			break
		}
		io.CopyN(newFile, conn, BUFFERSIZE)
		//Increment the counter
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}
