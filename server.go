package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type User struct {
	ID   int
	Conn *websocket.Conn
}

type Room struct {
	Name  string
	Users map[int]*User
	mu    sync.RWMutex
}

type WSMessage struct {
	To    int    `json:"to"`
	Event string `json:"event"`
	Data  string `json:"data"`
}

type RWSMessage struct {
	From  int    `json:"from"`
	Event string `json:"event"`
	Data  string `json:"data"`
}

type RoomData struct {
	Name string `json:"name"`
}

var (
	rooms   = make(map[string]*Room)
	roomsMu sync.RWMutex
)

func (r *Room) addUser(u *User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Users[u.ID] = u
}

func (r *Room) removeUser(id int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Users, id)
}

func (r *Room) getUserCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Users)
}

func reader(conn *websocket.Conn, room *Room, user *User) {
	defer func() {
		conn.Close()
		room.removeUser(user.ID)
		log.Printf("User %d left room %s. Room now has %d users.\n", user.ID, room.Name, room.getUserCount())
		if room.getUserCount() == 0 {
			roomsMu.Lock()
			delete(rooms, room.Name)
			roomsMu.Unlock()
		}
	}()

	for {
		var wsmsg WSMessage
		err := conn.ReadJSON(&wsmsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading WebSocket message: %v", err)
			}
			break
		}

		if wsmsg.Event == "createAnswer" || wsmsg.Event == "addIceCandidate" || wsmsg.Event == "createOffer" {
			rwsmsg := RWSMessage{Event: wsmsg.Event, Data: wsmsg.Data, From: user.ID}
			room.mu.RLock()
			targetUser, exists := room.Users[wsmsg.To]
			room.mu.RUnlock()
			if exists {
				err := targetUser.Conn.WriteJSON(rwsmsg)
				if err != nil {
					log.Printf("Error writing to WebSocket: %v", err)
				}
			}
		}

		log.Printf("Room: %s, Event: %s, From: %d, To: %d", room.Name, wsmsg.Event, user.ID, wsmsg.To)
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("roomname")

	roomsMu.RLock()
	room, exists := rooms[roomName]
	roomsMu.RUnlock()

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	id := room.getUserCount()
	user := &User{ID: id, Conn: ws}
	room.addUser(user)

	log.Printf("New client connected to room %s. Total users: %d", room.Name, room.getUserCount())

	reader(ws, room, user)
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	var rdata RoomData
	if err := json.NewDecoder(r.Body).Decode(&rdata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	roomsMu.Lock()
	if _, exists := rooms[rdata.Name]; exists {
		roomsMu.Unlock()
		http.Error(w, "Room already exists", http.StatusConflict)
		return
	}
	rooms[rdata.Name] = &Room{Name: rdata.Name, Users: make(map[int]*User)}
	roomsMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Room created successfully"})
}

func getRoomSize(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("roomname")

	roomsMu.RLock()
	room, exists := rooms[roomName]
	roomsMu.RUnlock()

	if !exists {
		w.Write([]byte("-1"))
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"size": room.getUserCount()})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./client/public")))
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/createRoom", createRoom)
	http.HandleFunc("/getRoomSize", getRoomSize)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
