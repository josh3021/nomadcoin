package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

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

func persistBlock(b *Block) {
	dbStorage.SaveBlock(b.Hash, utils.ToBytes(b))
}

func createBlock(previousHash string, height, difficulty int) *Block {
	block := &Block{
		Hash:         "",
		PreviousHash: previousHash,
		Height:       height,
		Difficulty:   difficulty,
		Nonce:        0,
	}
	block.Transactions = Mempool().ConfirmTxs()
	block.mine()
	fmt.Printf("\nHeight: %d\nHash: %s\nDifficulty: %d\nNonce: %d\n\n", block.Height, block.Hash, block.Difficulty, block.Nonce)
	persistBlock(block)
	return block
}

// FindBlock finds and returns block in database
func FindBlock(hash string) (*Block, error) {
	blockBytes := dbStorage.FindBlock(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}
