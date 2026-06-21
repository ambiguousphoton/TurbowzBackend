package main

import (
	"GoServer/authenticator"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
	UserMux.Lock()
	UsersList[userID] = conn
	log.Printf("User %s registered, total users: %d", userID, len(UsersList))
	UserMux.Unlock()
	broadcastPresence(userID, "online")
	deliverPendingMessages(userID)
}

func unregisterUser(userID string) {
	UserMux.Lock()
	delete(UsersList, userID)
	UserMux.Unlock()
	broadcastPresence(userID, "offline")
}

func broadcastPresence(userID, status string) {
	msg := map[string]string{
		"type":   "presence",
		"userID": userID,
		"status": status,
	}
	UserMux.Lock()
	defer UserMux.Unlock()
	for uid, conn := range UsersList {
		if uid != userID {
			conn.WriteJSON(msg)
		}
	}
}

func isUserOnline(userID string) bool {
	UserMux.Lock()
	defer UserMux.Unlock()
	_, exists := UsersList[userID]
	return exists
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


func sendMessageToRoom(sourceID, roomID, messageText, links, messageID, mediaURL, mediaType string) {
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
            "mediaURL":      mediaURL,
            "mediaType":     mediaType,
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
	log.Printf("Message sent to %d users in room %s with messageID %s", sentCount, roomID, messageID)
}




func handleNewMessage(message map[string]string){
	messageID     := uuid.New().String()
	sourceID      := message["sourceID"]
	destinationID := message["destinationID"]
	roomID        := message["roomID"]
	messageText   := message["messageText"]
	links         := message["links"]
	mediaURL      := message["mediaURL"]
	mediaType     := message["mediaType"] // "image", "video", "gif"

	log.Printf("Handling message from %s to %s in room %s", sourceID, destinationID, roomID)
	if destinationID != "" {
		roomID = MakeDuoRoom(sourceID, destinationID)
		log.Printf("Created new duo room %s for users %s and %s", roomID, sourceID, destinationID)
	}
	
	sendMessageToRoom(sourceID, roomID, messageText, links, messageID, mediaURL, mediaType)
	printQueue()
}


 
func handleReadReceipt(message map[string]string) {
	readerID := message["sourceID"]
	messageID := message["messageID"]
	senderID := message["senderID"]

	receipt := map[string]string{
		"type":      "read_receipt",
		"messageID": messageID,
		"readerID":  readerID,
	}

	UserMux.Lock()
	conn, online := UsersList[senderID]
	UserMux.Unlock()

	if online {
		conn.WriteJSON(receipt)
	}
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

		switch message["type"] {
		case "read_receipt":
			handleReadReceipt(message)
		default:
			handleNewMessage(message)
		}
	}

	unregisterUser(connecterIDstr)
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

func onlineStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "missing userID", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"online": isUserOnline(userID)})
}

const mediaDir = "MediaData/chat_media"

func chatMediaUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(50 << 20) // 50MB max

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	fileID := uuid.New().String() + ext

	os.MkdirAll(mediaDir, 0755)
	dst, err := os.Create(filepath.Join(mediaDir, fileID))
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	mediaURL := fmt.Sprintf("/chat-media?id=%s", fileID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mediaURL": mediaURL})
}

func chatMediaServeHandler(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join(mediaDir, filepath.Base(fileID))
	http.ServeFile(w, r, filePath)
}

func main(){
	http.HandleFunc("/connect-with-socket-server", authenticator.RequireAuth(websocketHandler))
	http.HandleFunc("/online-status", onlineStatusHandler)
	http.HandleFunc("/upload-chat-media", authenticator.RequireAuth(chatMediaUploadHandler))
	http.HandleFunc("/chat-media", chatMediaServeHandler)
	log.Println("Starting CommunicationService on Port 8280")
	err := http.ListenAndServe(":8280", nil)
	log.Println(err)
}