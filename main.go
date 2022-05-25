package main

import (
	"github.com/josh3021/nomadcoin/cli"
	"github.com/josh3021/nomadcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
	// wallet.Wallet()
}
