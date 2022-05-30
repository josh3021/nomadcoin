package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/josh3021/nomadcoin/utils"
)

var upgrader = websocket.Upgrader{}

// Upgrade upgrades http to ws protocol
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	queries := r.URL.Query()
	openPort := queries.Get("openPort")
	fmt.Println(openPort)
	utils.HandleErr(err)
	ipAddress := utils.Splitter(r.RemoteAddr, ":", 0)
	fmt.Println(openPort)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ipAddress != "" && openPort != ""
	}
	initPeer(conn, ipAddress, openPort)
}

// AddPeer adds a peer
func AddPeer(address string, port, openPort string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	initPeer(conn, address, port)
}
