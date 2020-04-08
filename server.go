package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type User struct {
	Id   int
	Conn *websocket.Conn
}

type Room struct {
	Name  string
	Users map[int]*User
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

var rooms = make(map[string]*Room, 5)

func reader(conn *websocket.Conn, room *Room, user *User) {
	ip := conn.RemoteAddr()
	var wsmsg WSMessage
	var rwsmsg RWSMessage

	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			log.Println("error occurred while reading websocket message")
			log.Println(err)
			continue
		}

		log.Println(string(m))

		err = json.Unmarshal(m, &wsmsg)
		if err != nil {
			log.Println("error occurred while unmarshalling websocket message")
			log.Panicln(err)
			continue
		}

		if wsmsg.Event == "createAnswer" || wsmsg.Event == "addIceCandidate" || wsmsg.Event == "createOffer" {
			rwsmsg = RWSMessage{Event: wsmsg.Event, Data: wsmsg.Data, From: user.Id}
			room.Users[wsmsg.To].Conn.WriteJSON(rwsmsg)
		}

		log.Printf("%s %s %s", ip, wsmsg.Event, wsmsg.Data)
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	roomname := queryValues.Get("roomname")
	room := rooms[roomname]
	if len(room.Name) > 0 {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		log.Println("client connected")

		id := 0
		l := len(room.Users)
		if l != 0 {
			id = room.Users[l-1].Id + 1
		}

		user := &User{Id: id, Conn: ws}
		room.Users[id] = user

		ws.SetCloseHandler(func(code int, text string) error {
			delete(room.Users, id)
			log.Printf("someone leaving room %s \nthis room has %d users...\n", room.Name, len(room.Users))
			if len(room.Users) == 0 {
				delete(rooms, roomname)
			}
			return ws.WriteMessage(websocket.TextMessage, []byte(text))
		})

		reader(ws, room, user)
	}
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("hola seniores")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body")
		log.Panicln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	var rdata RoomData
	err = json.Unmarshal(body, &rdata)
	rooms[rdata.Name] = &Room{Name: rdata.Name, Users: make(map[int]*User)}
}

func getRoomSize(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	roomname := queryValues.Get("roomname")
	room := rooms[roomname]
	log.Println("room: ", room)
	if len(roomname) == 0 || room == nil {
		w.Write([]byte("-1"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(strconv.Itoa(len(room.Users))))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./client/public")))
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/createRoom", createRoom)
	http.HandleFunc("/getRoomSize", getRoomSize)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
