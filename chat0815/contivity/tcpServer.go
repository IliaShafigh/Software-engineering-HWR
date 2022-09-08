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

func RunServer(l net.Listener, cStatusC chan *ChatroomStatus, refresh chan bool, errorC chan ErrorMessage) {
	log.Println("Listener initiating with server address", l.Addr().String())
	log.Println("SERVER: listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("SERVER: Error accepting incoming transmission from", conn.RemoteAddr().String())
			errorC <- ErrorMessage{Err: err, Msg: "Failed connection attempt from" + conn.RemoteAddr().String()}
		} else {
			log.Println("SERVER: Incoming TCP Request from", conn.RemoteAddr().String())
			go HandleRequest(conn, cStatusC, refresh, errorC)
		}
	}
}

//TODO if unknownIP(conn.Addr) && request != "UXXX {
//			perform UXXX
//		}
func HandleRequest(conn net.Conn, cStatusC chan *ChatroomStatus, refresh chan bool, errorC chan ErrorMessage) {
	log.Println("SERVER: TCP Accepted from", conn.RemoteAddr().String(), ",reading request type now...")
	//Expecting request type
	tmp := make([]byte, 70)
	_, err := conn.Read(tmp)
	request := string(tmp)[0:4]
	if err != nil {
		log.Println("SERVER: Could not Read request type because of:", err)
		return
	}
	log.Println("SERVER: Received request type " + request + "!")
	log.Println("SERVER: Full Message:", string(tmp))

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
			errorC <- ErrorMessage{Err: err, Msg: "SERVER: Could not encode and send cStatus"}
		}
		cStatusC <- cStatus
		log.Println("SERVER: Encoding is over!")

	case request == "GBXX":
		log.Println("SERVER: someone said goodbye, deleting", conn.RemoteAddr().String())
		RemoveUserAddr(conn.RemoteAddr(), cStatusC)
	case request == "NVKR":
	case request == "NGR":
	}

}

func RemoveUserAddr(toRemove net.Addr, cStatusC chan *ChatroomStatus) {
	cStatus := <-cStatusC
	for i, usrAddr := range cStatus.UserAddr {
		if strings.Split(toRemove.String(), ":")[0] == strings.Split(usrAddr.String(), ":")[0] {
			//Addr found so remove it and append everything else
			part2 := cStatus.UserAddr[i+1:]
			cStatus.UserAddr = cStatus.UserAddr[0:i]
			cStatus.UserAddr = append(cStatus.UserAddr, part2...)
			cStatusC <- cStatus
			log.Println(cStatus.UserAddr, "Removed ", toRemove.String())
			return
		}

	}
	cStatusC <- cStatus
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

type ErrorMessage struct {
	Err error
	Msg string
}
