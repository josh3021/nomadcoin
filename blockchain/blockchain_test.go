package blockchain

import (
	"reflect"
	"sync"
	"testing"

	"github.com/josh3021/nomadcoin/utils"
)

type fakeDB struct {
	fakeFindBlock      func() []byte
	fakeLoadBlockChain func() []byte
}

func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock()
}
func (f fakeDB) LoadBlockchain() []byte {
	return f.fakeLoadBlockChain()
}
func (fakeDB) SaveBlock(hash string, data []byte) {}
func (fakeDB) SaveBlockchain(data []byte)         {}
func (fakeDB) DeleteAllBlocks()                   {}

func TestBlockChain(t *testing.T) {
	t.Run("Should create Blockchain", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeLoadBlockChain: func() []byte {
				return nil
			},
		}
		bc := Blockchain()
		if bc.Height != 1 {
			t.Error("Blockchain() should create blockchain.")
		}
	})
	t.Run("Should restore Blockchain", func(t *testing.T) {
		once = *new(sync.Once)
		dbStorage = fakeDB{
			fakeLoadBlockChain: func() []byte {
				bc := &blockchain{Height: 2, CurrentDifficulty: 1, NewestHash: "XXXX"}
				return utils.ToBytes(bc)
			},
		}
		bc := Blockchain()
		if bc.Height != 2 {
			t.Errorf("Blockchain() should restore a blockchain with a height of %d, got %d", 2, bc.Height)
		}
	})
}

func TestBlocks(t *testing.T) {
	blocks := []*Block{
		{PreviousHash: "x"},
		{PreviousHash: ""},
	}
	fakeBlock := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func() {
				fakeBlock++
			}()
			return utils.ToBytes(blocks[fakeBlock])
		},
	}
	bc := &blockchain{}
	blocksResult := Blocks(bc)
	if reflect.TypeOf(blocksResult) != reflect.TypeOf([]*Block{}) {
		t.Error("Blocks() should return a slice of blocks")
	}
}

func TestFindTx(t *testing.T) {
	t.Run("Tx not found.", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height:       2,
					Transactions: []*Tx{},
				}
				return utils.ToBytes(b)
			},
		}
		tx := FindTx(&blockchain{}, "test")
		if tx != nil {
			t.Error("Tx should not be found.")
		}
	})
	t.Run("Tx found.", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 2,
					Transactions: []*Tx{
						{ID: "test"},
					},
				}
				return utils.ToBytes(b)
			},
		}
		tx := FindTx(&blockchain{NewestHash: "test"}, "test")
		if tx == nil {
			t.Error("Tx should be found.")
		}
	})
}

func TestGetDifficulty(t *testing.T) {
	blocks := []*Block{
		{PreviousHash: "x"},
		{PreviousHash: "x"},
		{PreviousHash: "x"},
		{PreviousHash: "x"},
		{PreviousHash: ""},
	}
	fakeBlock := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func() {
				fakeBlock++
			}()
			return utils.ToBytes(blocks[fakeBlock])
		},
	}
	type test struct {
		height int
		want   int
	}
	tests := []test{
		{height: 0, want: defaultDifficulty},
		{height: 2, want: 1},
		{height: 5, want: 5},
	}
	for _, tc := range tests {
		bc := &blockchain{Height: tc.height, CurrentDifficulty: defaultDifficulty}
		got := getDifficulty(bc)
		if got != tc.want {
			t.Errorf("getDifficulty() should return %d got %d", tc.want, got)
		}
	}
}

func TestAddPerrBlock(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "test",
	}
	Mempool().Txs["test"] = &Tx{}
	newBlock := &Block{
		Difficulty: 2,
		Hash:       "test",
		Transactions: []*Tx{
			{ID: "test"},
		},
	}
	bc.AddPeerBlock(newBlock)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("AddPeerBlock should mutate blockchain")
	}
}

func TestReplace(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "test",
	}
	newBlocks := []*Block{
		{Height: 2, Hash: "yy", Difficulty: 2},
		{Height: 1, Hash: "xx", Difficulty: 1},
	}
	bc.Replace(newBlocks)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "yy" {
		t.Error("Replace should replace blockchain")
	}
}
