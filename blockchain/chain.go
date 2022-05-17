package blockchain

import (
	"sync"

	"github.com/josh3021/nomadcoin/db"
	"github.com/josh3021/nomadcoin/utils"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) recalculateDifficulty() int {
	blocks := Blockchain().Blocks()
	latestBlock := blocks[0]
	lastRecalculatedBlock := blocks[difficultyInterval-1]
	actualInterval := (latestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedInterval := difficultyInterval * blockInterval

	if actualInterval <= (expectedInterval - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualInterval >= (expectedInterval + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) Difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	} else {
		return Blockchain().CurrentDifficulty
	}
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PreviousHash != "" {
			hashCursor = block.PreviousHash
		} else {
			break
		}
	}
	return blocks
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}
