package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/josh3021/nomadcoin/blockchain"
	"github.com/josh3021/nomadcoin/utils"
)

var upgrader = websocket.Upgrader{}

// Upgrade upgrades http to ws protocol
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	queries := r.URL.Query()
	openPort := queries.Get("openPort")
	utils.HandleErr(err)
	ipAddress := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ipAddress != "" && openPort != ""
	}
	initPeer(conn, ipAddress, openPort)
	// time.Sleep(10 * time.Second)
	// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from %s", )))

}

// AddPeer adds a peer
func AddPeer(address string, port, openPort string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	sendNewestBlock(peer)
	// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from %s", port)))
}

func BroadcastNewMessage(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.V {
		notifyNewMessage(b, p)
	}
}
