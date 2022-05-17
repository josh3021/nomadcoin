package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/josh3021/nomadcoin/db"
	"github.com/josh3021/nomadcoin/utils"
)

const (
	prefixTarget string = "0"
)

var ErrNotFound = errors.New("block not found")

type Block struct {
	Data         string `json:"data"`
	Hash         string `json:"hash"`
	PreviousHash string `json:"previousHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *Block) mine() {
	prefix := strings.Repeat(prefixTarget, b.Difficulty)
	for {
		hash := utils.Hash(b)
		fmt.Printf("hash: %s\n, nonce: %d\n", hash, b.Nonce)
		if strings.HasPrefix(hash, prefix) {
			b.Timestamp = int(time.Now().Unix())
			b.Hash = hash
			break
		}
		b.Nonce += 1
	}
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func createBlock(data string, previousHash string, height int) *Block {
	block := &Block{
		Data:         data,
		Hash:         "",
		PreviousHash: previousHash,
		Height:       height,
		Difficulty:   Blockchain().Difficulty(),
		Nonce:        0,
	}
	block.mine()
	block.persist()
	return block
}
