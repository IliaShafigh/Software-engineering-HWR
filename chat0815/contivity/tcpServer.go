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
	UserNames   map[string]string //UserNames[AddrWithoutPort(net.Addr.String())]name
	UserName    string
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
		AddGroupMessage(msg, conn.RemoteAddr(), cStatusC)
		refresh <- true
	case request == "UXXX":
		log.Println("SERVER: new Update request, encoding now... ")
		name := strings.TrimPrefix(string(tmp), request+":")
		//Add Addr
		if AddUserAddr(conn.RemoteAddr(), name, cStatusC) {
			//TODO SEND NEW IP TO ALL CLIENTS
			defer GUXX(cStatusC)
		}
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
	case request == "GUXX":
		log.Println("SERVER: new Get Update request, requesting now...")
		addr := net.TCPAddr{
			IP:   net.ParseIP(AddrWithoutPort(conn.RemoteAddr())),
			Port: 8888,
			Zone: "",
		}
		err = UXXX(&addr, cStatusC, refresh, errorC)
		if err != nil {
			errorC <- ErrorMessage{Err: err, Msg: "SERVER: Could not Get Updates from" + addr.String()}
		}
	case request == "GBXX":
		log.Println("SERVER: someone said goodbye, deleting", conn.RemoteAddr().String())
		RemoveUserAddr(conn.RemoteAddr(), cStatusC)
	case request == "NPMX":
		//TODO NEW PRIVATE MESSAGE
	case request == "NGR":
	}

}

//Adds User IP and Name to CStatus.
//Returns false if the address was already added.
func AddUserAddr(newAddr net.Addr, name string, cStatusC chan *ChatroomStatus) bool {
	cStatus := <-cStatusC
	for _, usrAddr := range cStatus.UserAddr {
		if strings.Split(newAddr.String(), ":")[0] == strings.Split(usrAddr.String(), ":")[0] {
			//Addr is already in s.UserAddr so nothing happens
			cStatusC <- cStatus
			return false
		}
	}
	toAdd := net.TCPAddr{
		IP:   net.ParseIP(strings.Split(newAddr.String(), ":")[0]),
		Port: 8888,
		Zone: "",
	}
	cStatus.UserAddr = append(cStatus.UserAddr, &toAdd)

	cStatus.UserNames[AddrWithoutPort(&toAdd)] = name
	cStatusC <- cStatus
	return true
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

func AddGroupMessage(msg string, senderAddr net.Addr, cStatusC chan *ChatroomStatus) {
	cStatus := <-cStatusC
	msg = cStatus.UserNames[AddrWithoutPort(senderAddr)] + ": " + msg
	cStatus.ChatContent = append(cStatus.ChatContent, msg)
	cStatusC <- cStatus
}

func TcpAddr(ip net.IP) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   ip,
		Port: 8888,
		Zone: "",
	}
}

type ErrorMessage struct {
	Err error
	Msg string
}

func mergeCStatus(newStatus ChatroomStatus, senderAddr net.Addr, cStatusC chan *ChatroomStatus) ChatroomStatus {
	//ChatContent Merge
	// TODO Improve chat merge maybe with timestamps
	//
	ownStatus := <-cStatusC
	if len(ownStatus.ChatContent) >= len(newStatus.ChatContent) {
		for _, msg := range ownStatus.ChatContent {
			//Compare
			//Do nothing because we assume that our chat is more advanced and we have the same messages
			_ = msg
		}
	} else {
		for i := len(ownStatus.ChatContent); i < len(newStatus.ChatContent); i++ {
			newMsgs := newStatus.ChatContent
			ownStatus.ChatContent = append(ownStatus.ChatContent, newMsgs[i])
		}
	}
	//UserAddr Merge
	for _, nAddr := range newStatus.UserAddr {
		if !contains(ownStatus.UserAddr, nAddr) {
			ownStatus.UserAddr = append(ownStatus.UserAddr, nAddr)
		}
	}
	//BlockedAddr Merge
	for _, nAddr := range newStatus.BlockedAddr {
		if !contains(ownStatus.UserAddr, nAddr) {
			ownStatus.BlockedAddr = append(ownStatus.BlockedAddr, nAddr)
		}
	}
	ownStatus.UserNames[AddrWithoutPort(senderAddr)] = newStatus.UserName
	cStatusC <- ownStatus
	return *ownStatus
}

func AddrWithoutPort(addr net.Addr) string {
	return strings.Split(addr.String(), ":")[0]
}

func PrintCStatus(cStatus ChatroomStatus) {
	log.Println("ChatContent", cStatus.ChatContent)
	log.Println("UserAddr", cStatus.UserAddr)
	log.Println("BlockedAddr", cStatus.BlockedAddr)
	log.Println("UserNames", cStatus.UserNames)
	log.Println("UserName", cStatus.UserName)
}
