package main

import (
	"fmt"
	"time"

	"github.com/josh3021/nomadcoin/cli"
	"github.com/josh3021/nomadcoin/db"
)

func sendOnly(c chan<- int) {
	for i := range [15]int{} {
		fmt.Printf(">> Sending: %d <<\n", i)
		time.Sleep(time.Second * 3)
		c <- i
		fmt.Printf(">> Sent: %d <<\n", i)
	}
	close(c)
}

func receiveOnly(c <-chan int) {
	for {
		i, ok := <-c
		fmt.Printf(">> Received: %d <<\n", i)
		if !ok {
			break
		}
		// fmt.Println(i)
	}
}

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
	// wallet.Wallet()
	// c := make(chan int, 10)
	// go sendOnly(c)
	// receiveOnly(c)
}
