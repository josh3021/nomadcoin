package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data         string
	Hash         string
	PreviousHash string
}

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PreviousHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

type blockchain struct {
	blocks []*Block
}

func (bc *blockchain) GetAllBlocks() []*Block {
	return bc.blocks
}

func (bc *blockchain) AddBlock(data string) {
	bc.blocks = append(bc.blocks, createBlock(data))
}

var once sync.Once
var bc *blockchain

func GetBlockChain() *blockchain {
	if bc == nil {
		once.Do(func() {
			bc = &blockchain{}
			bc.blocks = append(bc.blocks, createBlock("Genesis"))
		})
	}
	return bc
}

func getLastHash() string {
	blocks := GetBlockChain().blocks
	bcLength := len(blocks)
	if bcLength == 0 {
		return ""
	}
	return blocks[bcLength-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}
