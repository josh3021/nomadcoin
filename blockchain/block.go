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

// ErrNotFound returns ERROR if block not found.
var ErrNotFound = errors.New("block not found")

// Block is struct of the block in the blockchain.
type Block struct {
	Height       int    `json:"height"`
	Hash         string `json:"hash"`
	PreviousHash string `json:"previousHash,omitempty"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
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
		if strings.HasPrefix(hash, prefix) {
			b.Timestamp = int(time.Now().Unix())
			b.Hash = hash
			break
		}
		b.Nonce++
	}
}

// FindBlock finds and returns block in database
func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func createBlock(previousHash string, height int) *Block {
	block := &Block{
		Hash:         "",
		PreviousHash: previousHash,
		Height:       height,
		Difficulty:   Blockchain().Difficulty(),
		Nonce:        0,
		Transactions: []*Tx{makeCoinbaseTx("me")},
	}
	block.mine()
	fmt.Printf("\nHeight: %d\nHash: %s\nDifficulty: %d\nNonce: %d\n\n", block.Height, block.Hash, block.Difficulty, block.Nonce)
	block.persist()
	return block
}
