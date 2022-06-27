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
	m                 sync.Mutex
}

type storage interface {
	FindBlock(hash string) []byte
	SaveBlock(hash string, data []byte)
	SaveBlockchain(data []byte)
	LoadBlockchain() []byte
	DeleteAllBlocks()
}

const (
	defaultDifficulty  int = 4
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

var b *blockchain
var once sync.Once
var dbStorage storage = db.DB{}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock() *Block {
	// b.m.Lock()
	// defer b.m.Unlock()
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
	return block
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	b.Height = len(newBlocks)
	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	dbStorage.DeleteAllBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height++
	b.CurrentDifficulty = newBlock.Difficulty
	b.NewestHash = newBlock.Hash

	persistBlockchain(b)
	persistBlock(newBlock)

	for _, tx := range newBlock.Transactions {
		_, ok := m.Txs[tx.ID]
		if ok {
			delete(m.Txs, tx.ID)
		}
	}
}

func persistBlockchain(b *blockchain) {
	dbStorage.SaveBlockchain(utils.ToBytes(b))
}

func recalculateDifficulty(b *blockchain) int {
	blocks := Blocks(b)
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

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return Blockchain().CurrentDifficulty
	}
}

// Blocks returns all blocks
func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()
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

// UTxOutsByAddress returns Unspent Transaction Outputs By Address
func UTxOutsByAddress(b *blockchain, address string) []*UTxOut {
	var uTxOuts []*UTxOut
	var creatorTxs = make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, txIn := range tx.TxIns {
				if txIn.Signature == "COINBASE" {
					break
				}
				if FindTx(b, txIn.TxID).TxOuts[txIn.Index].Address == address {
					creatorTxs[txIn.TxID] = true
				}
			}

			for index, txOut := range tx.TxOuts {
				if txOut.Address == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{
							TxID:   tx.ID,
							Index:  index,
							Amount: txOut.Amount,
						}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

// GetBalanceByAddress returns balance of address
func GetBalanceByAddress(b *blockchain, address string) int {
	var balance int
	ownedTxOuts := UTxOutsByAddress(b, address)
	for _, ownedTxOut := range ownedTxOuts {
		balance += ownedTxOut.Amount
	}
	return balance
}

// Txs return all transactions
func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

// FindTx returns tx that want
func FindTx(b *blockchain, targetTxID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetTxID {
			return tx
		}
	}
	return nil
}

func Status() *blockchain {
	b.m.Lock()
	defer b.m.Unlock()
	return b
}

// Blockchain returns blockchain (Initialize blockchain if it does not initialized).
func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := dbStorage.LoadBlockchain()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			b.restore(checkpoint)
		}
	})
	return b
}
