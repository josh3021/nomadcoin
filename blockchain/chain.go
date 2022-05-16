package blockchain

import (
	"errors"
	"sync"

	"github.com/josh3021/nomadcoin/db"
	"github.com/josh3021/nomadcoin/utils"
)

var NotFoundError = errors.New("Block Not Found.")

type blockchain struct {
	Height   int    `json:"height"`
	LastHash string `json:"lastHash"`
}

var once sync.Once
var bc *blockchain

func (bc *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(bc))
}

func (bc *blockchain) AddBlock(data string) {
	block := createBlock(data, bc.LastHash, bc.Height+1)
	bc.Height = block.Height
	bc.LastHash = block.Hash
	bc.persist()
}

func BlockChain() *blockchain {
	if bc == nil {
		once.Do(func() {
			bc = &blockchain{0, ""}
			bc.AddBlock("Genesis")
		})
	}
	return bc
}
