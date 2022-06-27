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
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)
	// time.Sleep(10 * time.Second)
	// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from %s", )))

}

// AddPeer adds a peer
func AddPeer(address, port, openPort string, isBroadcast bool) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	if isBroadcast {
		broadcastNewPeer(peer)
		return
	}
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

func BroadcastNewTx(tx *blockchain.Tx) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.V {
		notifyNewTx(tx, p)
	}
}

func broadcastNewPeer(newPeer *peer) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for key, p := range Peers.V {
		if newPeer.key != key {
			address := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(address, p)
		}
	}
}
