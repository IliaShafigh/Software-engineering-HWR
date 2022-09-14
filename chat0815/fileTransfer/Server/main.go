package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

// Define the size of how big the chunks of data will be send each time
// can be between 1 to 65495 bytes
const BUFFERSIZE = 1024

func main() {
	//Create a TCP listener on port 8888
	server, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Starting Server failed:", err)
		os.Exit(1)
	}
	defer server.Close() //case main ends
	fmt.Println("Server started! Waiting for connections...")
	//Spawn a new goroutine whenever a client connects
	for {
		conn, err := server.Accept() //acccepts coonection
		if err != nil {
			fmt.Println("Server could not established connection:", err)
			os.Exit(1)
		}
		fmt.Println("Client connected") //send file to specific target(conn)
		go sendFileToClient(conn)
	}
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

// function to send file

func sendFileToClient(connection net.Conn) {
	fmt.Println("A client has connected!")
	defer connection.Close()
	//file to sned t client
	file, err := os.Open("aFile.txt")
	if err != nil {
		fmt.Println("Couldn't open file:", err)
		return
	}
	defer file.Close()
	//Get the filename and filesize
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	//hier werden die fehlenden bytes von name und size aufgefüllt
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	//Send the file header first so the client knows the filename and how long it has to read the incomming file
	fmt.Println("Sending filename and filesize!")
	//Write first 10 bytes to client telling them the filesize
	connection.Write([]byte(fileSize))
	//Write 64 bytes to client containing the filename
	connection.Write([]byte(fileName))
	//Initialize a buffer for reading parts of the file in
	sendBuffer := make([]byte, BUFFERSIZE)
	//Start sending the file to the client
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			//End of file reached, break out of for loop
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return
}
