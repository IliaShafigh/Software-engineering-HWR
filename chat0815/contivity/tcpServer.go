package contivity

import (
	"log"
	"net"
)

type ChatroomStatus struct {
	ChatContent *[]string
	UserAddr    *[]net.TCPAddr
	BlockedAddr *[]net.TCPAddr
}
type ChatroomStatusTcp struct {
	ChatContent []string
	UserAddr    []net.TCPAddr
	BlockedAddr []net.TCPAddr
}

func packCStatus(cStatus ChatroomStatus) ChatroomStatusTcp {
	return ChatroomStatusTcp{
		ChatContent: *cStatus.ChatContent,
		UserAddr:    *cStatus.UserAddr,
		BlockedAddr: *cStatus.BlockedAddr,
	}

}

func RunServer(l net.Listener, cStatus ChatroomStatus, refresh chan bool) {
	log.Println("Listener initiating with server address", l.Addr().String())

	log.Println("SERVER: listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("SERVER: Error accepting incoming transmission from", conn.RemoteAddr().String())
			log.Fatal(err)
		}
		log.Println("SERVER: Incoming TCP Request from", conn.RemoteAddr().String(), ". Sending ChatRoomStatus...")

		go HandleRequest(conn, cStatus, refresh)
	}
}

func HandleRequest(conn net.Conn, cStatus ChatroomStatus, refresh chan bool) {
	log.Println("SERVER: TCP Accepted from", conn.RemoteAddr().String())
	//Expecting request type
	tmp := make([]byte, 1024)
	_, err := conn.Read(tmp)
	request := string(tmp)[0:4]

	if err != nil {
		log.Println("SERVER: Could not Read request type because of:", err)
	}
	log.Println("SERVER: Received request type " + request + "!")

	switch {
	case request == "NGMX":
		//Confirm request
		log.Println("SERVER: new group Message request")
		_, err = conn.Write([]byte("New Group Message X"))
		if err != nil {
			log.Println("SERVER: could not write confirmation for", request, "type")
			return
		}
		//Read New Group Message
		newGroupMsg := make([]byte, 1024)
		_, err = conn.Read(newGroupMsg)
		if err != nil {
			log.Println("SERVER: could not read new group Message from", conn.RemoteAddr().String())
			return
		}
		*cStatus.ChatContent = append(*cStatus.ChatContent, string(newGroupMsg))
		refresh <- true
	case request == "UXXX":
		//Confirm request
		log.Println("SERVER: new Update request ")
		//Confirm Update Request X
		_, err = conn.Write([]byte("CURX"))
		if err != nil {
			log.Println("SERVER: could not write confirmation for", request, "type")
		}

		err = sendChatContent(conn, cStatus)
		if err != nil {
			//TODO abort?
		}
		err = sendAddresses(conn, cStatus)
		if err != nil {

		}
		conn.Close()
	case request == "NVKR":
	case request == "NGR":
	}

}

func sendAddresses(conn net.Conn, cStatus ChatroomStatus) error {
	//Send User Adresses
	for _, addr := range *cStatus.UserAddr {
		_ = addr
	}
	//Send Blocked User Adresses
	return nil
}

func sendChatContent(conn net.Conn, cStatus ChatroomStatus) error {
	log.Println("SERVER: sending ChatContent")
	for _, msg := range *cStatus.ChatContent {
		//Write message
		log.Println("SERVER: trying to send...")
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Println("SERVER: could not write msg ", msg, "to", conn.RemoteAddr().String(), "SKIPPING")
		}
	}

	_, err := conn.Write([]byte("XXX"))
	if err != nil {
		log.Println("SERVER: Could not write end of ChatContent to", conn.RemoteAddr().String())
		return err
	}
	return nil
}
