package contivity

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
)

func GetStatusUpdate(addr net.Addr, cStatus ChatroomStatus) {
	//Connection
	log.Println("Client: Trying to connect to", addr.String())
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Println("Client: could not establish connection to get Updates from", addr.String())
		return
	}
	defer conn.Close()
	log.Println("Client: connected successfully to:", conn.RemoteAddr().String())

	//Indicate Request and wait for confirmation
	err = confirmRequest("UXXX", conn)
	if err != nil {
		log.Println("Client: could not confirm Update Request from", addr.String())
		return
	}

	//Updating cStatus
	readChatContent(conn, cStatus)
	readAddresses(conn, cStatus)

	//write Message to indicate exchange is over
	//TODO Maybe useless cause conn.Close()
	//_, err = conn.Write([]byte("XXX"))
	log.Println("Client: Got all updates, closing connection now")
	return
}

//TODO
//Send Message to all participants of the Group
func SendMessageToGroup(cStatus ChatroomStatus, msg string) {
	log.Println("Client: Sending Message to Group")
	request := "NGMX"
	for _, addr := range *cStatus.UserAddr {
		err := sendMsg(addr, msg, request)
		if err != nil {
			//TODO if not reachable, delete from cStatus?
			log.Println("Client: Could not send Group Message to:", addr.String(), ", SKIPPING")
			continue
		}
	}
	return
}

func sendMsg(addr net.TCPAddr, msg string, request string) error {
	//Connect to addr
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Println("Client: could not connect to:", conn.RemoteAddr())
		return err
	}
	defer conn.Close()
	//Confirm request
	err = confirmRequest(request, conn)
	if err != nil {
		log.Println("Client: could not Confirm:", request, "to:", addr.String())
		return err
	}
	//send new Message
	_, err = conn.Write([]byte(msg))
	if err != nil {
		log.Println("Client: could not send Message of request type", request, "to:", addr.String())
		return err
	}
	return nil
}

//
func confirmRequest(request string, conn net.Conn) error {
	//TODO: add more request types
	tmp := make([]byte, 1024)
	switch request {
	case "UXXX": //Confirm update request
		_, err := conn.Write([]byte(request))
		if err != nil {
			log.Println("Client: could not send request type:", request)
			return err
		}

		_, err = conn.Read(tmp)
		confirmation := string(tmp)[:4]
		if err != nil && err != io.EOF {
			log.Println("Client: could not read request type confirmation:", request, err)
			return err
		} else if string(confirmation) != "CURX" {
			log.Println("Client:", conn.RemoteAddr().String(), "did not confirm correctly to request type:", request, string(confirmation))
			return errors.New("Client: ELSE THEN CONFIRMATION")
		}
		log.Println("Client: confirmed request type:", request)
		return nil
	case "NGMX": //New Group Message X
		_, err := conn.Write([]byte(request))
		if err != nil {
			log.Println("Client: could not send request type:", request)
			return err
		}
		_, err = conn.Read(tmp)
		confirmation := string(tmp)[:4]
		if err != nil {
			log.Println("Client: could not read request type confirmation:", request)
			return err
		} else if string(confirmation) != "New Group Message X" {
			log.Println("Client:", conn.RemoteAddr().String(), "did not confirm request type:", request)
			return errors.New("Client: ELSE THEN CONFIRMATION")
		}
		log.Println("Client: confirmed request type:", request)
		return nil

	case "":

	}
	return errors.New("something went wrong CONFIRM REQUEST")
}

//TODO (errors?)
func readAddresses(conn net.Conn, status ChatroomStatus) {
	log.Println("Client: reading UserAddr...")

	log.Println("Client: UserAddr successfully read!")
	log.Println("Client: reading BlockedAddr...")

	log.Println("Client: BlockedAddr successfully read!")
}

//TODO return errors? or abort
func readChatContent(conn net.Conn, cStatus ChatroomStatus) {
	log.Println("Client: reading ChatContent...")

	for {
		//Read Message
		log.Println("Client: trying to read...")
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			continue
		}
		if msg[:3] == "XXX" {
			break
		}
		log.Println("Client: Received msg", msg)
		*cStatus.ChatContent = append(*cStatus.ChatContent, msg)
		log.Println("Client: Appended", msg)
	}
}
