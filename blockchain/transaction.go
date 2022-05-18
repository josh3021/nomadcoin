package blockchain

import (
	"errors"
	"time"

	"github.com/josh3021/nomadcoin/utils"
)

const minerReward int = 50

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (tx *Tx) getId() {
	tx.ID = utils.Hash(tx)
}

type TxIn struct {
	Owner  string
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
}

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("Not Enough Money.")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	total := 0
	for _, oldTxOut := range oldTxOuts {
		if amount <= total {
			break
		}
		oldTxOuts = append(oldTxOuts, oldTxOut)
		total += oldTxOut.Amount
	}
	change := total - amount

	if change > 0 {
		changeTxOut := &TxIn{Owner: from, Amount: change}
		txIns = append(txIns, changeTxOut)
	}
	txOut := &TxOut{Owner: to, Amount: amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("me", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{Owner: "COINBASE", Amount: minerReward},
	}
	txOuts := []*TxOut{
		{Owner: address, Amount: minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}
