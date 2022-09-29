package contivity

import (
	"encoding/gob"
	"log"
	"net"
	"strings"
)

// UXXX Get Status Update request. Name is equal to request switch on tcpServer.go
func UXXX(addr net.Addr, cStatusC chan *GroupChatStatus, refresh chan bool, errorC chan ErrorMessage) error {
	//Connection
	log.Println("Client: Trying to connect to", addr.String())
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Println("Client: could not establish connection to get Updates from", addr.String(), err)
		refresh <- false
		errorC <- ErrorMessage{Err: err, Msg: "Client: could not establish connection to get Updates from, " + addr.String()}
		return err
	}
	defer conn.Close()
	log.Println("Client: connected successfully to:", conn.RemoteAddr().String())
	log.Println("Client: writing request type with own name...")
	tmp := <-cStatusC
	cStatusC <- tmp
	name := tmp.UserName
	_, err = conn.Write([]byte("UXXX:" + name))
	if err != nil {
		log.Println("Client: could not write request type cause of:", err)
		errorC <- ErrorMessage{Err: err, Msg: "Client: could not write request type to, " + addr.String()}
	}
	log.Println("Client: did write request type, trying to decode now...")

	decoder := gob.NewDecoder(conn)
	gob.Register(&net.TCPAddr{})
	newCStatus := &GroupChatStatus{}
	err = decoder.Decode(newCStatus)
	if err != nil {
		log.Println("Client: Problem with Decoding cause of", err)
		errorC <- ErrorMessage{Err: err, Msg: "Client: could not decode cStatus from, " + addr.String()}
	}
	log.Println("Client: Deconding seems to have worked")
	log.Println("Client: this is the Status from Remote:")
	PrintCStatus(*newCStatus)
	log.Println("Client: this is the own Status")
	tmp = <-cStatusC
	cStatusC <- tmp
	cStatus := *tmp
	PrintCStatus(cStatus)
	cStatus = mergeCStatus(*newCStatus, conn.RemoteAddr(), cStatusC)
	log.Println("Client: this is the merged Status")
	PrintCStatus(cStatus)
	refresh <- true
	log.Println("Client: Got all updates, closing connection now")
	return nil
}

// GUXX sends Get Update Request to all Participants of own cStatus.
//Request that all receivers send UXXX request to oneself.
//Name is equal to request switch on tcpServer.go
func GUXX(cStatusC chan *GroupChatStatus) {
	log.Println("Client: sending GUXX Request to everybody now...")
	request := "GUXX"

	cStatus := <-cStatusC
	userAdresses := cStatus.UserAddr
	cStatusC <- cStatus
	for _, addr := range userAdresses {
		if TcpAddr(GetOutboundIP()).String() != addr.String() {
			err := sendMsg(addr, "", request)
			if err != nil {
				//TODO if not reachable, delete from cStatus?
				log.Println("Client: Could not send Group Message to:", addr.String(), ", SKIPPING")
				continue
			}
		}
	}
}

// NGMX sends Message to all group members. send to all participants of the Group including oneself
//Updates of ChatDisplay should be implemented in tcpServer
//Name is equal to request switch on tcpServer.go
func NGMX(msg string, cStatusC chan *GroupChatStatus, errorC chan ErrorMessage) {
	log.Println("Client: Sending Message to Group")
	//errorC <- ErrorMessage{Err: nil, Msg: "Client msg: " + msg + " "}
	request := "NGMX"

	tmp := <-cStatusC
	userAddresses := tmp.UserAddr
	cStatusC <- tmp
	for _, addr := range userAddresses {
		err := sendMsg(addr, msg, request)
		if err != nil {
			//TODO if not reachable, delete from cStatus?
			log.Println("Client: Could not send Group Message to:", addr.String(), ", SKIPPING")
			continue
		}
	}
	return
}

// GBXX Say Goodbye to all participants in your cStatus
//Name is equal to request switch on tcpServer.go
func GBXX(cStatusC chan *GroupChatStatus) {
	log.Println("Saying goodbye now!")
	request := "GBXX"
	tmp := <-cStatusC
	userAddresses := tmp.UserAddr
	cStatusC <- tmp
	for _, addr := range userAddresses {
		err := sendMsg(addr, "", request)
		if err != nil {
			//TODO if not reachable, delete from cStatus?
			log.Println("Client: Could not send Group Message to:", addr.String(), ", SKIPPING")
			continue
		}
	}
}
func contains(addrs []net.Addr, addr2 net.Addr) bool {
	if addr2 == nil {
		return false
	}
	for _, addr := range addrs {

		if addr.String() == addr2.String() {
			return true
		}
	}
	return false
}

func SendMessageToPrivate(msg string, addr net.Addr, errorC chan ErrorMessage) {
	log.Println("Client: sending Private message to", addr.String())
	request := "NPMX"
	err := sendMsg(addr, msg, request)
	if err != nil {
		errorC <- ErrorMessage{Err: err, Msg: "Could not send message to" + addr.String() + " | " + msg}
		log.Println("Client: failed to send Private Message to", addr.String())
	}

}

func sendMsg(addr net.Addr, msg string, request string) error {
	//Connect to addr
	connectAddr := strings.Split(addr.String(), ":")[0] + ":8888"
	conn, err := net.Dial("tcp", connectAddr)
	if err != nil {
		log.Println("Client: conn err :", err, addr.String())
		return err
	}
	defer conn.Close()
	//Write request type and new msg
	_, err = conn.Write([]byte(request + ":" + msg))
	if err != nil {
		log.Println("Client: could not write request type cause of:", err)
		return err
	}
	log.Println("Client: did write request type and send msg to", addr.String())
	return nil
}

// GetOutboundIP we need own IP for many reasons
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
