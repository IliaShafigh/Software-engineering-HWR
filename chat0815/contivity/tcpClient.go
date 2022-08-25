package contivity

import (
	"log"
	"net"
)

func ConnectPLS(addr net.Addr, done chan bool) {
	connection, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Println(1)
		log.Fatal(err)
	}
	defer connection.Close()
	msg := "Penis"
	_, err = connection.Write([]byte(msg))
	if err != nil {
		log.Println(2)
		log.Fatal(err)
	}
	reply := make([]byte, 1024)
	_, err = connection.Read(reply)
	if err != nil {
		log.Println(3)
		log.Fatal(err)
	}
	log.Println("I am the user and i received:", string(reply))
	done <- true

}

//TODO
//Send Message to all participants of the Group
func SendMessageToGroup(message string) {

}
