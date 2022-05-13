package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Height       int    `json:"height"`
	Data         string `json:"data"`
	Hash         string `json:"hash"`
	PreviousHash string `json:"previousHash,omitempty"`
}

var NotFoundError = errors.New("Block Not Found.")

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

func (bc *blockchain) GetBlock(height int) (*Block, error) {
	if height > len(bc.GetAllBlocks()) {
		return nil, NotFoundError
	}
	return bc.blocks[height-1], nil
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
	newBlock := Block{
		Height:       len(GetBlockChain().GetAllBlocks()) + 1,
		Data:         data,
		Hash:         "",
		PreviousHash: getLastHash(),
	}
	newBlock.calculateHash()
	return &newBlock
}
