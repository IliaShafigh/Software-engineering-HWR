package contivity

import (
	"chat0815/errPopUps"
	"chat0815/fileTransfer"
	"encoding/gob"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"log"
	"net"
	"strings"
)

type ChatStorage struct {
	*container.AppTabs // corresponding apptabs in which our chats tabitems are stored
	Navigation         *fyne.Container
	*GroupChat
	Private    []*PrivateChat
	MainWindow fyne.Window
}
type GroupChat struct {
	*container.TabItem                 // TabItem with display and entry
	Navigation         *fyne.Container // the left side of our window, so called navigation
	GcStatusC          chan *GroupChatStatus
	Refresh            chan bool
}
type GroupChatStatus struct {
	ChatContent []string
	UserAddr    []net.Addr
	BlockedAddr []net.Addr
	UserNames   map[string]string //UserNames[AddrWithoutPort(net.Addr.String())]name
	UserName    string            //OWN NAME
}

type PrivateChat struct {
	*container.TabItem // TabItem with display and Entry
	PvStatusC          chan *PrivateChatStatus
	Refresh            chan bool       //Refreshes Display
	Navigation         *fyne.Container //should include buttons for Hung, Hai und Ilia
}

type PrivateChatStatus struct {
	ChatContent []string
	UserAddr    net.Addr //Addr from remote partner of the private Chat
	Ttg         *TTGGameStatus
	Sv          *SVGameStatus
}

type TTGGameStatus struct {
	Running   bool
	Won       bool
	MyTurn    bool
	GameField [9]int
}

type SVGameStatus struct {
	MyTurn  bool
	Running bool
	Won     bool
}

func RunServer(l net.Listener, chatC chan ChatStorage, errorC chan errPopUps.ErrorMessage) {
	log.Println("Listener initiating with server address", l.Addr().String())
	log.Println("SERVER: listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("SERVER: Error accepting incoming transmission ", err)
			errorC <- errPopUps.ErrorMessage{Err: err, Msg: "Failed connection attempt "}
		} else {
			log.Println("SERVER: Incoming TCP Request from", conn.RemoteAddr().String())
			go HandleRequest(conn, chatC, errorC)
		}
	}
}

//TODO if unknownIP(conn.Addr) && request != "UXXX {
//			perform UXXX
//		}
func HandleRequest(conn net.Conn, chatC chan ChatStorage, errorC chan errPopUps.ErrorMessage) {
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
	chats := <-chatC
	gcStatusC := chats.GroupChat.GcStatusC
	refresh := chats.GroupChat.Refresh
	chatC <- chats

	switch {
	case request == "NGMX": // New Group Message
		log.Println("SERVER: new Group Message requets")
		msg := strings.TrimPrefix(string(tmp), request+":")
		log.Println("SERVER: msg received was:", msg)
		AddGroupMessage(msg, conn.RemoteAddr(), gcStatusC)
		refresh <- true
	case request == "UXXX": //Update Request
		log.Println("SERVER: new Update request, encoding now... ")
		name := strings.TrimPrefix(string(tmp), request+":")
		//Add Addr
		if AddUserAddr(conn.RemoteAddr(), name, gcStatusC) {
			//TODO SEND NEW IP TO ALL CLIENTS
			defer GUXX(gcStatusC)
		}
		gcStatus := <-gcStatusC
		encoder := gob.NewEncoder(conn)
		gob.Register(&net.TCPAddr{})
		err = encoder.Encode(*gcStatus)
		if err != nil {
			log.Println("SERVER: Problem with encoding:", err)
			errorC <- errPopUps.ErrorMessage{Err: err, Msg: "SERVER: Could not encode and send gcStatus"}
		}
		gcStatusC <- gcStatus
		log.Println("SERVER: Encoding is over!")
	case request == "GUXX": //Get Update Request
		log.Println("SERVER: new Get Update request, requesting now...")
		addr := net.TCPAddr{
			IP:   net.ParseIP(AddrWithoutPort(conn.RemoteAddr())),
			Port: 8888,
			Zone: "",
		}
		err = UXXX(&addr, chatC, refresh, errorC)
		if err != nil {
			errorC <- errPopUps.ErrorMessage{Err: err, Msg: "SERVER: Could not Get Updates from" + addr.String()}
		}
	case request == "GBXX": //Good Bye Request
		log.Println("SERVER: someone said goodbye, deleting", conn.RemoteAddr().String())
		RemoveUserAddr(conn.RemoteAddr(), chatC)
	case request == "NPMX": // New Private Message
		log.Println("SERVER: new Private Message requets")
		msg := strings.TrimPrefix(string(tmp), request+":")
		log.Println("SERVER: Private msg received was:", msg)
		chats = <-chatC
		for indexOCPT, pvChat := range chats.Private {
			pvStatus := <-pvChat.PvStatusC
			pvChat.PvStatusC <- pvStatus
			if AddrWithoutPort(conn.RemoteAddr()) == AddrWithoutPort(pvStatus.UserAddr) {
				gcStatus := <-chats.GcStatusC
				name := gcStatus.UserNames[AddrWithoutPort(pvStatus.UserAddr)]
				chats.GcStatusC <- gcStatus
				chatC <- chats
				//Private Chat Tab is already open
				AddPrivateMessage(msg, chatC, indexOCPT, name)
				return
			}
		}
		//TODO Got a new private Message so open new tab
		//if its the second entry we will run here even though there is a tab existing. so export this
		chatC <- chats
		//TODO Refresh Display of private Chat?
	case request == "NFTX": //New File Transfer Request
		log.Println("SERVER: new File Transfer request")
		//TODO Open Private Chat First or give some kind of notification or accept stuff
		fileTransfer.SaveFile(conn, chats.MainWindow, errorC)
	}

}

//Adds User IP and Name to CStatus.
//Returns false if the address was already added.
func AddUserAddr(newAddr net.Addr, name string, gcStatusC chan *GroupChatStatus) bool {
	gcStatus := <-gcStatusC
	for _, usrAddr := range gcStatus.UserAddr {
		if strings.Split(newAddr.String(), ":")[0] == strings.Split(usrAddr.String(), ":")[0] {
			//Addr is already in s.UserAddr so nothing happens
			gcStatusC <- gcStatus
			return false
		}
	}
	toAdd := net.TCPAddr{
		IP:   net.ParseIP(strings.Split(newAddr.String(), ":")[0]),
		Port: 8888,
		Zone: "",
	}
	gcStatus.UserAddr = append(gcStatus.UserAddr, &toAdd)

	gcStatus.UserNames[AddrWithoutPort(&toAdd)] = name
	gcStatusC <- gcStatus
	return true
}

func RemoveUserAddr(toRemove net.Addr, chatC chan ChatStorage) {
	chats := <-chatC
	gcStatus := <-chats.GcStatusC
	for i, usrAddr := range gcStatus.UserAddr {
		if strings.Split(toRemove.String(), ":")[0] == strings.Split(usrAddr.String(), ":")[0] {
			//Addr found so remove it and append everything else
			part2 := gcStatus.UserAddr[i+1:]
			gcStatus.UserAddr = gcStatus.UserAddr[0:i]
			gcStatus.UserAddr = append(gcStatus.UserAddr, part2...)
			log.Println(gcStatus.UserAddr, "Removed ", toRemove.String())
			break
		}
	}
	delete(gcStatus.UserNames, AddrWithoutPort(toRemove))
	PrintCStatus(*gcStatus)
	chats.Navigation.Remove(chats.Navigation.Objects[0])
	chats.Navigation.Add(chats.GroupChat.Navigation)
	chats.GcStatusC <- gcStatus
	chatC <- chats

}

func AddGroupMessage(msg string, senderAddr net.Addr, gcStatusC chan *GroupChatStatus) {
	gcStatus := <-gcStatusC
	msg = gcStatus.UserNames[AddrWithoutPort(senderAddr)] + ": " + msg
	gcStatus.ChatContent = append(gcStatus.ChatContent, msg)
	gcStatusC <- gcStatus
}

func AddPrivateMessage(msg string, chatC chan ChatStorage, indexOCPT int, prefixName string) {
	chats := <-chatC
	gcStatus := <-chats.GroupChat.GcStatusC
	pvStatus := <-chats.Private[indexOCPT].PvStatusC
	msg = prefixName + ": " + msg
	pvStatus.ChatContent = append(pvStatus.ChatContent, msg)
	chats.Private[indexOCPT].Refresh <- true
	chats.Private[indexOCPT].PvStatusC <- pvStatus
	chats.GroupChat.GcStatusC <- gcStatus
	chatC <- chats
}

// TcpAddr Takes ip and adds port 8888 and returns net.TCPAddr
func TcpAddr(ip net.IP) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   ip,
		Port: 8888,
		Zone: "",
	}
}

func mergeCStatus(newStatus GroupChatStatus, senderAddr net.Addr, gcStatusC chan *GroupChatStatus) GroupChatStatus {
	//ChatContent Merge
	// TODO Improve chat merge maybe with timestamps
	//
	ownStatus := <-gcStatusC
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
	for nUserNameKey, nUserName := range newStatus.UserNames {
		_, exists := ownStatus.UserNames[nUserNameKey]
		if !exists {
			//Add username and address
			ownStatus.UserNames[nUserNameKey] = nUserName
		}
	}
	gcStatusC <- ownStatus
	return *ownStatus
}

func AddrWithoutPort(addr net.Addr) string {
	return strings.Split(addr.String(), ":")[0]
}

func PrintCStatus(gcStatus GroupChatStatus) {
	log.Println("ChatContent", gcStatus.ChatContent)
	log.Println("UserAddr", gcStatus.UserAddr)
	log.Println("BlockedAddr", gcStatus.BlockedAddr)
	log.Println("UserNames", gcStatus.UserNames)
	log.Println("UserName", gcStatus.UserName)
}

// InitializeGroupChatRoomStatus Should only be called once for initialization
func InitializeGroupChatRoomStatus() *GroupChatStatus {
	chatContent := make([]string, 0)

	chatContent = append(chatContent, "Take care of each other and watch your drink")
	chatContent = append(chatContent, "Welcome to chat0815")

	gcStatus := GroupChatStatus{
		ChatContent: chatContent,
		UserAddr:    []net.Addr{},
		BlockedAddr: []net.Addr{},
		UserNames:   make(map[string]string), //map[IPADRESSE]Name
		UserName:    "",
	}
	//Fill own information
	gcStatus.UserAddr = append(gcStatus.UserAddr, TcpAddr(GetOutboundIP()))
	return &gcStatus
}

// InitializePrivateChatRoomStatus Should only be called once for initialization
func InitializePrivateChatRoomStatus(remoteAddr net.Addr) *PrivateChatStatus {
	chatContent := make([]string, 0)

	chatContent = append(chatContent, "This is private Chat")

	pvStatus := PrivateChatStatus{
		ChatContent: chatContent,
		UserAddr:    remoteAddr,
		Ttg: &TTGGameStatus{
			Running:   false,
			Won:       false,
			MyTurn:    false,
			GameField: [9]int{},
		},
		Sv: &SVGameStatus{
			MyTurn:  false,
			Running: false,
			Won:     false,
		},
	}
	return &pvStatus
}
