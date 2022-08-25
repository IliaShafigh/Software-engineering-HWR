package contivity

import (
	"log"
	"net"
)

func RunServer(l net.Listener) {

	for {
		friend, err := l.Accept()
		if err != nil {
			log.Println("Error accepting")
			log.Fatal(err)
		}
		log.Println("Message from", friend.RemoteAddr(), "to", friend.LocalAddr())
		go HandleRequest(friend)
	}
}

func HandleRequest(friend net.Conn) {
	buf := make([]byte, 1024)
	_, err := friend.Read(buf)
	if err != nil {
		log.Println("Could not Read Message because of:", err)
	}
	if string(buf) == "ende" {
		panic("Ende befehl gelesen")
	}
	friend.Write([]byte("Thank you for writing:" + string(buf) + " to me you a******"))
	friend.Close()

}
