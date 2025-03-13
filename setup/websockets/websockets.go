package websockets

import (
	"go-server/setup"
	"go-server/setup/reqRes"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (customize for security)
	},
}

type Peer struct {
	Id   string
	Conn *websocket.Conn
}

var (
	peers   = make(map[string]*Peer)
	peersMu sync.Mutex
)

func parseMessage(connection *websocket.Conn) (map[string]any, error) {
	var msg map[string]any
	err := connection.ReadJSON(&msg)
	return msg, err
}

func rememberPeer(peer *Peer) {
	peersMu.Lock()
	defer peersMu.Unlock()
	peers[peer.Id] = peer
}

func forgetPeer(peer *Peer) {
	peersMu.Lock()
	defer peersMu.Unlock()
	delete(peers, peer.Id)
}

type MessageAssignId struct {
	YourId string `json:"yourId"`
}

type MessageUserLeaves struct {
	FromId string `json:"fromId"`
	Bye    bool   `json:"bye"`
}

func HandleWebSocket(w reqRes.MyWriter, r *reqRes.MyRequest) {
	conn, err := upgrader.Upgrade(w, &r.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error error: \n%v\n", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("WebSocket connection close error: \n%v\n", err)
		}
	}(conn)

	peer := &Peer{Id: setup.RandId(), Conn: conn}

	messageAssignId := MessageAssignId{
		YourId: peer.Id,
	}

	ok := sendMessage(peer, messageAssignId)
	if !ok {
		log.Printf("Failed to send messageAssignId to peer %s, I guess we just drop the connection.\n", peer.Id)
		return
	}

	log.Printf("New WebRTC peer connected: %s\n", peer.Id)
	rememberPeer(peer)

	for {
		msg, err := parseMessage(conn)
		if err != nil {
			log.Printf("Error reading JSON: \n%v\n", err)
			break
		}

		handleMessage(peer.Id, msg)
	}

	// tell everyone that they are going away
	handleGlobalMessage(peer.Id, MessageUserLeaves{FromId: peer.Id, Bye: true})

	// Cleanup on disconnect
	forgetPeer(peer)
	log.Printf("WebRTC peer disconnected: %s\n", peer.Id)
}

func sendMessage(peer *Peer, message any) bool {
	err := peer.Conn.WriteJSON(message)
	if err != nil {
		log.Printf("WebSocket send error to peer %s: \n%v\n", peer.Id, err)
		return false
	}
	return true
}

func handleMessage(senderId string, message map[string]any) {
	targetId, hasTarget := message["toId"].(string)
	if hasTarget {
		handleTargetedMessage(senderId, targetId, message)
	} else {
		handleGlobalMessage(senderId, message)
	}
}

func handleTargetedMessage(senderId string, targetId string, message map[string]any) {
	peersMu.Lock()
	defer peersMu.Unlock()

	recipient, found := peers[targetId]
	if !found {
		log.Printf("Sender %s tried to send to target %s but they were not found!\n", senderId, targetId)
		return
	}

	err := recipient.Conn.WriteJSON(message)
	if err != nil {
		log.Printf("Sender %s tried to send to target %s but sending failed: \n%v\n", senderId, targetId, err)
	}
}

func handleGlobalMessage(senderId string, message any) {
	peersMu.Lock()
	defer peersMu.Unlock()

	for _, peer := range peers {
		if peer.Id == senderId {
			continue
		}

		err := peer.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("Sender %s tried to send to target %s but sending failed: \n%v\n", senderId, peer.Id, err)
		}
	}
}
