package db

import (
	"fmt"
	"os"

	"github.com/josh3021/nomadcoin/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName       = "blockchain"
	dbExtName    = ".db"
	dataBucket   = "data"
	blocksBucket = "blocks"

	checkpoint = "checkpoint"
)

var db *bolt.DB

func getDBName() string {
	port := os.Args[2][10:]
	return fmt.Sprintf("%s_%s%s", dbName, port, dbExtName)
}

// DB returns database (Initialize database if it does not initialized).
func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDBName(), 0600, nil)
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

// Close Database.
func Close() {
	DB().Close()
}

// SaveBlock saves the block in database.
func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

// SaveBlockchain saves the blockchain in database.
func SaveBlockchain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

// Checkpoint returns the checkpoint from database
func Checkpoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

// FindBlock returns the block from database
func FindBlock(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

// EmptyBlocks delete and recreate blocksBucket
func EmptyBlocks() {
	DB().Update(func(tx *bolt.Tx) error {
		utils.HandleErr(tx.DeleteBucket([]byte(blocksBucket)))
		_, err := tx.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
