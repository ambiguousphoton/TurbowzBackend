package main

import (
	"GoServer/authenticator"
	"encoding/json"
	"log"
	"net/http"
	// "os/user"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	// "golang.org/x/tools/go/analysis/passes/defers"
)


var (
    PendingMessages = make(map[string][]map[string]string) // userID -> list of messages
    PendingMux       sync.Mutex
)



var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}


var (
	UsersList = make(map[string]*websocket.Conn)  // map of userID to their websocket connection
	UserMux sync.Mutex                // mutex to protect concurrent access to the
	RoomsList = make(map[string]map[string]struct{})   // {"roomid1": {"vy333":{}, "Omi":{}}, "roomid3": {"vy334":{}, "Omi4":{}}}
	RoomsMux sync.Mutex
)



func registerUser(userID string, conn *websocket.Conn){
	UserMux.Lock();
	UsersList[userID] = conn;
	log.Printf("User %s registered, total users: %d", userID, len(UsersList))
	UserMux.Unlock();
	deliverPendingMessages(userID)
}




func MakeRoom(userIDs []string) string {
	roomID := uuid.New().String()
	RoomsMux.Lock()
	defer RoomsMux.Unlock()

	if _, exists := RoomsList[roomID]; !exists {
		RoomsList[roomID] = make(map[string]struct{})
	}
	for _, uid := range userIDs {
		RoomsList[roomID][uid] = struct{}{}
	}
	log.Printf("Room %s created with users: %v", roomID, userIDs)
	return roomID
}


func MakeDuoRoom(user1ID, user2ID string) string {
	var roomID string
	if user1ID < user2ID {
		roomID = user1ID + "_" + user2ID
	} else {
		roomID = user2ID + "_" + user1ID
	}
	
	RoomsMux.Lock()
	defer RoomsMux.Unlock()

	if _, exists := RoomsList[roomID]; !exists {
		RoomsList[roomID] = make(map[string]struct{})
		RoomsList[roomID][user1ID] = struct{}{}
		RoomsList[roomID][user2ID] = struct{}{}
		log.Printf("Room %s created with users: %v", roomID, []string{user1ID, user2ID})
	}
	return roomID
}


func sendMessageToRoom(sourceID, roomID, messageText, links, messageID string) {
	RoomsMux.Lock()
	members, roomExists := RoomsList[roomID]
	RoomsMux.Unlock()

	if !roomExists {
		log.Println("Room does not exist: ", roomID)
		return
	}

	log.Printf("Sending message from %s to room %s: %s", sourceID, roomID, messageText)
	UserMux.Lock()
	defer UserMux.Unlock()
	sentCount := 0
	msg := map[string]string{
			"messageID":    messageID,
            "sourceID":      sourceID,
            "messageText":   messageText,
            "roomID":        roomID,
            "links":         links,
        }
	for userID := range members {
		if userID != sourceID {
			if userConn, exists := UsersList[userID]; exists {
				msg["destinationID"] = userID
				err := userConn.WriteJSON(msg)
				if err != nil {
					log.Printf("Failed to send message to user %s: %v", userID, err)
					delete(UsersList, userID)
                	queueMessage(userID, msg)
					log.Printf("User %s removed from active connections", userID)
				} else {
					sentCount++
				}
			} else {
				queueMessage(userID, msg)
				log.Printf("User %s not found in active connections", userID)
			}
		}
	}
	log.Printf("Message sent to %d users in room %s with messageID", sentCount, roomID, messageID)
}




func handleNewMessage(message map[string]string){
	messageID     := uuid.New().String()
	sourceID      := message["sourceID"]
	destinationID := message["destinationID"]
	roomID        := message["roomID"]
	messageText   := message["messageText"]
	links         := message["links"]

	log.Printf("Handling message from %s to %s in room %s", sourceID, destinationID, roomID)
	if destinationID != "" {
		roomID = MakeDuoRoom(sourceID, destinationID)
		log.Printf("Created new duo room %s for users %s and %s", roomID, sourceID, destinationID)
	}
	
	sendMessageToRoom(sourceID, roomID, messageText, links, messageID)
	printQueue()
}


 
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	connecterID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		http.Error(w, "Error Invalid UserId ", http.StatusBadRequest)
		log.Println("error Invalid UserId ", connecterID)
		return
	}
	connecterIDstr := strconv.FormatInt(connecterID, 10)
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	registerUser(connecterIDstr, conn)
	log.Println("User connected with ID: ", connecterID)

	for {
		var rawMessage []byte
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error from user %s: %v", connecterIDstr, err)
			break
		}
		
		var message map[string]string
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			log.Printf("JSON parse error from user %s: %v, raw message: %s", connecterIDstr, err, string(rawMessage))
			continue
		}
		
		message["sourceID"] = connecterIDstr
		handleNewMessage(message)
	}
}




func printQueue() {
    PendingMux.Lock()
    defer PendingMux.Unlock()
    
    log.Printf("=== PENDING MESSAGES QUEUE ===")
    for userID, messages := range PendingMessages {
        log.Printf("User %s: %d pending messages", userID, len(messages))
    }
    log.Printf("=== END QUEUE ===")
}

func queueMessage(userID string, message map[string]string) {
    PendingMux.Lock()
    defer PendingMux.Unlock()

    PendingMessages[userID] = append(PendingMessages[userID], message)
    log.Printf("Queued message for user %s", userID)
}

func deliverPendingMessages(userID string) {
    PendingMux.Lock()
    messages, hasPending := PendingMessages[userID]
    if hasPending {
        delete(PendingMessages, userID)
    }
    PendingMux.Unlock()

    if !hasPending {
        return
    }

    log.Printf("Delivering %d pending messages to %s", len(messages), userID)
    UserMux.Lock()
    conn, online := UsersList[userID]
    UserMux.Unlock()

    if !online {
        log.Printf("User %s disconnected before delivery", userID)
        // requeue the messages
        PendingMux.Lock()
        PendingMessages[userID] = append(PendingMessages[userID], messages...)
        PendingMux.Unlock()
        return
    }

    for i, msg := range messages {
        err := conn.WriteJSON(msg)
        if err != nil {
            log.Printf("Error delivering pending message to %s: %v", userID, err)
            // requeue remaining messages
            queueMessage(userID, msg)
        } else {
            log.Printf("Successfully delivered pending message %d to %s: %s", i+1, userID, msg["messageText"])
        }
    }
    log.Printf("Finished delivering pending messages to %s", userID)
}

func main(){
	http.HandleFunc("/connect-with-socket-server", authenticator.RequireAuth(websocketHandler));
	log.Println("Starting CommunicationService on Port 8280")
	err := http.ListenAndServe(":8280", nil)
	log.Println(err)
}