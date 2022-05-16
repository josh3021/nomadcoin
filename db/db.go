package db

import (
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/josh3021/nomadcoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
)

var once sync.Once
var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		once.Do(func() {
			dbP, err := bolt.Open(dbName, 0600, nil)
			utils.HandleErr(err)
			db = dbP
			err = db.Update(func(tx *bolt.Tx) error {
				_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
				utils.HandleErr(err)
				_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
				utils.HandleErr(err)
				return nil
			})
			utils.HandleErr(err)
		})
	}
	return db
}

func SaveBlock(hash string, data []byte) {
	fmt.Printf("Hash: %s\nData: %b\n", hash, data)
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockchain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte("checkpoint"), data)
		return err
	})
	utils.HandleErr(err)
}
