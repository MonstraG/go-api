package websockets

import (
	"errors"
	"go-api/helpers"
	"go-api/setup"
	"go-api/setup/myLog"
	"go-api/setup/reqRes"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // maybe add origins, but then I have to deal with nginx
	},
}

type Room struct {
	Id    string
	Peers map[string]*Peer
	Mutex sync.Mutex
}

type Peer struct {
	Id   string
	Room string
	Conn *websocket.Conn
}

var (
	rooms      = make(map[string]*Room)
	roomsMutex sync.Mutex
)

func getOrCreateRoom(roomId string) *Room {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, exists := rooms[roomId]
	if !exists {
		room = &Room{
			Id:    roomId,
			Peers: make(map[string]*Peer),
		}
		rooms[roomId] = room
	}
	return room
}

func getRoom(roomId string) (*Room, error) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	room, exists := rooms[roomId]
	if !exists {
		return nil, errors.New("Failed to find room " + roomId)
	}
	return room, nil
}

func parseMessage(connection *websocket.Conn) (map[string]any, error) {
	var msg map[string]any
	err := connection.ReadJSON(&msg)
	return msg, err
}

func rememberPeer(peer *Peer, room *Room) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.Peers[peer.Id] = peer
}

func forgetPeer(peer *Peer, room *Room) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	delete(room.Peers, peer.Id)

	// If the room is empty, remove it
	if len(room.Peers) == 0 {
		roomsMutex.Lock()
		defer roomsMutex.Unlock()
		delete(rooms, room.Id)
	}
}

type MessageAssignId struct {
	YourId string `json:"yourId"`
}

type MessageUserLeaves struct {
	FromId string `json:"fromId"`
	Bye    bool   `json:"bye"`
}

func HandleWebSocket(w reqRes.MyWriter, r *reqRes.MyRequest) {
	connection, err := upgrader.Upgrade(w, &r.Request, nil)
	if err != nil {
		myLog.Info.Logf("WebSocket upgrade error error: \n%v\n", err)
		return
	}
	defer helpers.CloseSafely(connection)

	peerId := setup.RandId()
	myLog.Info.Logf("New peer %s trying to connect, wating for the room message\n", peerId)

	initialMsg, err := parseMessage(connection)
	if err != nil {
		myLog.Info.Logf("Failed to read initial message: \n%v\n", err)
		return
	}

	roomId, ok := initialMsg["roomId"].(string)
	if !ok || roomId == "" {
		myLog.Info.Log("Client did not send a valid roomId, closing connection.")
		err := connection.Close()
		if err != nil {
			myLog.Info.Logf("WebSocket connection close error: \n%v\n", err)
		}
		return
	}

	room := getOrCreateRoom(roomId)
	peer := &Peer{Id: peerId, Room: room.Id, Conn: connection}

	messageAssignId := MessageAssignId{
		YourId: peer.Id,
	}

	ok = sendMessage(peer, messageAssignId)
	if !ok {
		myLog.Info.Logf("Failed to send messageAssignId to peer %s in room %s, I guess we just drop the connection.\n", peer.Id, peer.Room)
		return
	}

	myLog.Info.Logf("New WebRTC peer connected: %s in room %s\n", peer.Id, peer.Room)
	rememberPeer(peer, room)

	for {
		msg, err := parseMessage(connection)
		if err != nil {
			myLog.Info.Logf("Error reading JSON: \n%v\n", err)
			break
		}

		handleMessage(peer, msg)
	}

	// tell everyone that they are going away
	handleRoomMessage(peer, MessageUserLeaves{FromId: peer.Id, Bye: true})

	// Cleanup on disconnect
	forgetPeer(peer, room)
	myLog.Info.Logf("WebRTC peer disconnected: %s\n", peer.Id)
}

func sendMessage(peer *Peer, message any) bool {
	err := peer.Conn.WriteJSON(message)
	if err != nil {
		myLog.Info.Logf("WebSocket send error to peer %s: \n%v\n", peer.Id, err)
		return false
	}
	return true
}

func handleMessage(sender *Peer, message map[string]any) {
	targetId, hasTarget := message["toId"].(string)
	if hasTarget {
		handleTargetedMessage(sender, targetId, message)
	} else {
		handleRoomMessage(sender, message)
	}
}

func handleTargetedMessage(sender *Peer, targetId string, message map[string]any) {
	room, err := getRoom(sender.Room)
	if err != nil {
		myLog.Info.Logf("Sender %s tried to send to target %s but room %s was not found!\n", sender.Id, targetId, sender.Room)
		return
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	recipient, found := room.Peers[targetId]
	if !found {
		myLog.Info.Logf("Sender %s tried to send to target %s but they were not found!\n", sender.Id, targetId)
		return
	}

	err = recipient.Conn.WriteJSON(message)
	if err != nil {
		myLog.Info.Logf("Sender %s tried to send to target %s but sending failed: \n%v\n", sender.Id, targetId, err)
	}
}

func handleRoomMessage(sender *Peer, message any) {
	room, err := getRoom(sender.Room)
	if err != nil {
		myLog.Info.Logf("Sender %s tried to send to room %s, but it was not found!\n", sender.Id, sender.Room)
		return
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	for _, peer := range room.Peers {
		if peer.Id == sender.Id {
			continue
		}

		err = peer.Conn.WriteJSON(message)
		if err != nil {
			myLog.Info.Logf("Sender %s tried to send to target %s but sending failed: \n%v\n", sender.Id, peer.Id, err)
		}
	}
}
