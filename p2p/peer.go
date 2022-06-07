package p2p

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/josh3021/nomadcoin/utils"
)

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

type peers struct {
	V map[string]*peer
	m sync.Mutex
}

var Peers peers = peers{
	V: make(map[string]*peer),
}

func (p *peer) close() {
	Peers.m.Lock()
	defer func() {
		time.Sleep(time.Second * 20)
		Peers.m.Unlock()
	}()
	err := p.conn.Close()
	utils.HandleErr(err)
	delete(Peers.V, p.key)
}

func (p *peer) read() {
	defer p.close()
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m)
		if err != nil {
			break
		}
		handleMessage(&m, p)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
	// utils.HandleErr(err)
}

func AllPeers(p *peers) []string {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	var keys []string
	for key := range p.V {
		keys = append(keys, key)
	}
	return keys
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		key:     key,
		address: address,
		port:    port,
		conn:    conn,
		inbox:   make(chan []byte),
	}
	go p.read()
	go p.write()
	Peers.V[key] = p
	return p
}
