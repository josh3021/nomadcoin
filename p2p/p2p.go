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
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	for {
		_, p, err := ws.ReadMessage()
		utils.HandleErr(err)
		message := fmt.Sprintf("new message: %s\n", p)
		// fmt.Printf("new message: %s\n", message)
		utils.HandleErr(ws.WriteMessage(websocket.TextMessage, []byte(message)))
	}
}
