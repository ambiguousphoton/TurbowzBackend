package main

import (
	"GoServer/authenticator"
	"log"
	"net/http"
	// "os/user"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	// "golang.org/x/tools/go/analysis/passes/defers"
)






var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}


var (
	UsersList = make(map[string]*websocket.Conn)  // map of userID to their websocket connection
	UserMux sync.Mutex                // mutex to protect concurrent access to the

	RoomsList = make(map[string]map[string]struct{})
	RoomsMux sync.Mutex
)



func registerUser(userID string, conn *websocket.Conn){
	UserMux.Lock();
	defer UserMux.Unlock()	;
	UsersList[userID] = conn;
}




func MakeRoom(userIDs []string) string {
	roomID := uuid.New().String()
	RoomsMux.Lock()
	defer RoomsMux.Unlock()

	if _, RoomExists := RoomsList[roomID]; !RoomExists {
		RoomsList[roomID] = make(map[string]struct{})
	}
	for _, uid := range userIDs {
        RoomsList[roomID][uid] = struct{}{} // add user
    }

	return roomID
}


func MakeDuoRoom(user1ID, user2ID string) string {
	return MakeRoom([]string{user1ID, user2ID})
}


func sendMessageToRoom(roomID, senderID, text string) {
	RoomsMux.Lock()
	members, roomExists := RoomsList[roomID]
	RoomsMux.Unlock()

	if !roomExists {
		log.Println("Room does not exist: ", roomID)
		return
	}

	UserMux.Lock()
	defer UserMux.Unlock()
	for userID := range members {
		if userID != senderID {
			if userConn, exists := UsersList[userID]; exists {
				userConn.WriteJSON(map[string]string{
					"Type":     "Message",
					"SenderID": senderID,
					"Text":     text,
					"RoomID":   roomID,
				})
			}
		}
	}
}




func handleChat(senderID string, message map[string]string){
	receiverID := message["ReceiverID"]
	

	switch message["Type"] {
		case "Connection-Request":
			if user, userPresent := UsersList[receiverID]; userPresent{
				user.WriteJSON(map[string]string{
					"Type":      "Connection-Request",
					"SenderID":   senderID,
				})
			}

		case "Connection-Approval":
			status := message["Status"]
			if status == "Accepted"{
				DuoRoomID := MakeDuoRoom(senderID, receiverID)
				UserMux.Lock()
				for _, uid := range []string{senderID, receiverID} {
					otherUserID := receiverID
					if uid == receiverID {
						otherUserID = senderID
					}
					if userConn, exists := UsersList[uid]; exists {
						userConn.WriteJSON(map[string]string{
							"Type": "Start-Chat",
							"RoomID": DuoRoomID,
							"With": otherUserID,
						})
						log.Println("Connection approved, room created: ", DuoRoomID, " between ", senderID, " and ", receiverID)
					}else{
						if uid == ""{
							log.Println("No User Id: ", uid)
							continue
						}else {
							log.Println("User not found: ", uid)
						}
					}
				}
				UserMux.Unlock()
			}

		case "Message":
			roomID := message["RoomID"]
			text := message["Text"]
			sendMessageToRoom(roomID, senderID, text)
			


		}
	

}



func websocketHandler(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		http.Error(w, "Error Invalid UserId ", http.StatusBadRequest)
		log.Println("error Invalid UserId ", requesterID)
		return
	}
	requesterIDstr := strconv.FormatInt(requesterID, 10)
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	registerUser(requesterIDstr, conn)
	log.Println("User connected with ID: ", requesterID)

	for {
		var message map[string]string;
		if err := conn.ReadJSON(&message) ; err != nil {
			log.Println("Read error:", err)
			break
		}
		handleChat(requesterIDstr, message)
	}
}


func main(){
	http.HandleFunc("/chat", authenticator.RequireAuth(websocketHandler));
	log.Println("Starting SocketConnectionService on Port 8181")
	http.ListenAndServe(":8181", nil)
}