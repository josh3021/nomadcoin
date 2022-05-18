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

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
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

func (b *blockchain) txOuts() []*TxOut {
	var txOuts []*TxOut
	blocks := b.Blocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

func (b *blockchain) TxOutsByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	txOuts := b.txOuts()
	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut)
		}
	}
	return ownedTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	var balance int
	ownedTxOuts := b.TxOutsByAddress(address)
	for _, ownedTxOut := range ownedTxOuts {
		balance += ownedTxOut.Amount
	}
	return balance
}

// Blockchain returns blockchain (Initialize blockchain if it does not initialized).
func Blockchain() *blockchain { //nolint:golint
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock()
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}
