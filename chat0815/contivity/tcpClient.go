package contivity

import (
	"encoding/gob"
	"log"
	"net"
	"strings"
)

func GetStatusUpdate(addr net.Addr, cStatusC chan *ChatroomStatus, refresh chan bool) {
	//Connection
	log.Println("Client: Trying to connect to", addr.String())
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Println("Client: could not establish connection to get Updates from", addr.String())
		return
	}
	defer conn.Close()
	log.Println("Client: connected successfully to:", conn.RemoteAddr().String())
	log.Println("Client: writing request type...")
	_, err = conn.Write([]byte("UXXX"))
	if err != nil {
		log.Println("Client: could not write request type cause of:", err)
	}
	log.Println("Client: did write request type, trying to decode now...")

	decoder := gob.NewDecoder(conn)
	newCStatus := &ChatroomStatus{}
	err = decoder.Decode(newCStatus)
	if err != nil {
		log.Println("Client: Problem with Decoding cause of", err)
	}
	log.Println("Client: Deconding seems to have worked")
	log.Println("Client: this is the Status from Remote:")
	log.Println("Client:", newCStatus.ChatContent)
	log.Println("Client:", newCStatus.UserAddr)
	log.Println("Client:", newCStatus.BlockedAddr)
	log.Println("Client: this is the own Status")
	tmp := <-cStatusC
	cStatusC <- tmp
	cStatus := *tmp
	log.Println("Client:", cStatus.ChatContent)
	log.Println("Client:", cStatus.UserAddr)
	log.Println("Client:", cStatus.BlockedAddr)
	cStatus = mergeCStatus(*newCStatus, cStatusC)
	log.Println("Client: this is the merged Status")
	log.Println("Client:", cStatus.ChatContent)
	log.Println("Client:", cStatus.UserAddr)
	log.Println("Client:", cStatus.BlockedAddr)
	refresh <- true
	log.Println("Client: Got all updates, closing connection now")
	return
}

//TODO Improve maybe?
func mergeCStatus(newStatus ChatroomStatus, cStatusC chan *ChatroomStatus) ChatroomStatus {
	//ChatContent Merge
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
	if len(ownStatus.UserAddr) >= len(newStatus.UserAddr) {
		for _, msg := range ownStatus.UserAddr {
			//Compare
			//Do nothing because we assume that our chat is more advanced and we have the same messages
			_ = msg
		}
	} else {
		for i := len(ownStatus.UserAddr); i < len(newStatus.UserAddr); i++ {
			newAddrs := newStatus.UserAddr
			ownStatus.UserAddr = append(ownStatus.UserAddr, newAddrs[i])
		}
	}
	//BlockedAddr Merge
	if len(ownStatus.BlockedAddr) >= len(newStatus.BlockedAddr) {
		for _, msg := range ownStatus.BlockedAddr {
			//Compare
			//Do nothing because we assume that our chat is more advanced and we have the same messages
			_ = msg
		}
	} else {
		for i := len(ownStatus.BlockedAddr); i < len(newStatus.BlockedAddr); i++ {
			newBlAddrs := newStatus.BlockedAddr
			ownStatus.BlockedAddr = append(ownStatus.BlockedAddr, newBlAddrs[i])
		}
	}
	cStatusC <- ownStatus
	return *ownStatus
}

//Send Message to all participants of the Group including oneself
//Updates of ChatDisplay should be implemented in tcpServer
func SendMessageToGroup(msg string, cStatusC chan *ChatroomStatus) {
	log.Println("Client: Sending Message to Group")
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

func sendMsg(addr net.Addr, msg string, request string) error {
	//Connect to addr
	connectAddr := strings.Split(addr.String(), ":")[0] + ":8888"
	conn, err := net.Dial("tcp", connectAddr)
	if err != nil {
		log.Println("Client: could not connect to:", conn.RemoteAddr())
		return err
	}
	defer conn.Close()
	//Write request type and new msg
	_, err = conn.Write([]byte(request + ":" + msg))
	if err != nil {
		log.Println("Client: could not write request type cause of:", err)
	}
	log.Println("Client: did write request type and send msg to", addr.String())
	return nil
}

//Need own IP for many reasons
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
