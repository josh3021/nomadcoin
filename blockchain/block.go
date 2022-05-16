package blockchain

import (
	"crypto/sha256"
	"fmt"

	"github.com/josh3021/nomadcoin/db"
	"github.com/josh3021/nomadcoin/utils"
)

type Block struct {
	Height       int    `json:"height"`
	Data         string `json:"data"`
	Hash         string `json:"hash"`
	PreviousHash string `json:"previousHash,omitempty"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func createBlock(data string, previousHash string, height int) *Block {
	block := &Block{
		Data:         data,
		Hash:         "",
		PreviousHash: previousHash,
		Height:       height,
	}
	payload := block.Data + block.PreviousHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}
