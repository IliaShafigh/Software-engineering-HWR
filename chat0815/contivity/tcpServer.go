package contivity

import (
	"encoding/gob"
	"log"
	"net"
	"strings"
)

type ChatroomStatus struct {
	ChatContent []string
	UserAddr    []net.Addr
	BlockedAddr []net.Addr
}

func AddUserAddr(newAddr net.Addr, cStatusC chan *ChatroomStatus) {
	cStatus := <-cStatusC

	for _, usrAddr := range cStatus.UserAddr {
		if strings.Split(newAddr.String(), ":")[0] == strings.Split(usrAddr.String(), ":")[0] {
			//Addr is already in s.UserAddr so nothing happens
			cStatusC <- cStatus
			return
		}

	}
	toAdd := net.TCPAddr{
		IP:   net.ParseIP(strings.Split(newAddr.String(), ":")[0]),
		Port: 8888,
		Zone: "",
	}
	cStatus.UserAddr = append(cStatus.UserAddr, &toAdd)
	cStatusC <- cStatus
}

func RunServer(l net.Listener, cStatusC chan *ChatroomStatus, refresh chan bool) {
	log.Println("Listener initiating with server address", l.Addr().String())

	log.Println("SERVER: listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("SERVER: Error accepting incoming transmission from", conn.RemoteAddr().String())
			log.Fatal(err)
		}
		log.Println("SERVER: Incoming TCP Request from", conn.RemoteAddr().String())

		go HandleRequest(conn, cStatusC, refresh)
	}
}

func HandleRequest(conn net.Conn, cStatusC chan *ChatroomStatus, refresh chan bool) {
	log.Println("SERVER: TCP Accepted from", conn.RemoteAddr().String(), ",reading request type now...")
	//Expecting request type
	tmp := make([]byte, 8)
	_, err := conn.Read(tmp)
	request := string(tmp)[0:4]
	if err != nil {
		log.Println("SERVER: Could not Read request type because of:", err)
	}
	log.Println("SERVER: Received request type " + request + "!")

	switch {
	case request == "NGMX":
		log.Println("SERVER: new Group Message requets")
		msg := strings.TrimPrefix(string(tmp), request+":")
		log.Println("SERVER: msg received was:", msg)
		AddGroupMessage(msg, cStatusC)
		refresh <- true
	case request == "UXXX":
		log.Println("SERVER: new Update request, encoding now... ")
		//Add Addr
		AddUserAddr(conn.RemoteAddr(), cStatusC)

		cStatus := <-cStatusC
		encoder := gob.NewEncoder(conn)
		gob.Register(&net.TCPAddr{})
		err = encoder.Encode(*cStatus)
		if err != nil {
			log.Println("SERVER: Problem with encoding:", err)
			panic("")
		}
		cStatusC <- cStatus
		log.Println("SERVER: Encoding is over!")

	case request == "NVKR":
	case request == "NGR":
	}

}

func AddGroupMessage(msg string, cStatusC chan *ChatroomStatus) {
	cStatus := <-cStatusC
	cStatus.ChatContent = append(cStatus.ChatContent, msg)
	cStatusC <- cStatus
}

func TcpAddr(ip net.IP) net.TCPAddr {
	return net.TCPAddr{
		IP:   ip,
		Port: 8888,
		Zone: "",
	}
}
