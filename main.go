package main

import (
	"fmt"

	"github.com/josh3021/nomadcoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockChain()
	chain.AddBlock("Second")
	chain.AddBlock("Third")
	chain.AddBlock("Fourth")
	for index, block := range chain.GetAllBlocks() {
		fmt.Printf("%d's Block\n", index+1)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Previous Hash: %s\n", block.PrevHash)
	}
}
