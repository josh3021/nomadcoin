package blockchain

import (
	"reflect"
	"testing"

	"github.com/josh3021/nomadcoin/utils"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instance of a block")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Block should not be found.", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				return nil
			},
		}
		_, err := FindBlock("x")
		if err == nil {
			t.Error("Block should not be found.")
		}
	})
	t.Run("Block should be found.", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{}
				return utils.ToBytes(b)
			},
		}
		b, _ := FindBlock("x")
		if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
			t.Error("Block should be found.")
		}
	})
}
