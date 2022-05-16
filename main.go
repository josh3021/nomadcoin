package main

import (
	"github.com/josh3021/nomadcoin/blockchain"
)

func main() {
	blockchain.BlockChain().AddBlock("First")
	blockchain.BlockChain().AddBlock("Second")
	blockchain.BlockChain().AddBlock("Third")
}
