package blockchain

import (
	"errors"
	"time"

	"github.com/josh3021/nomadcoin/utils"
)

const minerReward int = 50

// Tx contains information of transactions
type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (tx *Tx) getID() {
	tx.ID = utils.Hash(tx)
}

// TxIn contains information of transactions input
type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

// TxOut contains information of transactions Output
type TxOut struct {
	Owner  string
	Amount int
}

// UTxOut contains information of Unconfirmed transactions Output
type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

type mempool struct {
	Txs []*Tx
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("me", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) ConfirmTxs() []*Tx {
	coinbase := makeCoinbaseTx("me")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}

// Mempool contains not confirmed transactions
var Mempool *mempool = &mempool{}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, txIn := range tx.TxIns {
			if txIn.TxID == uTxOut.TxID && txIn.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if GetBalanceByAddress(Blockchain(), from) < amount {
		return nil, errors.New("not enough money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	uTxOuts := UTxOutsByAddress(Blockchain(), from)
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{TxID: uTxOut.TxID, Index: uTxOut.Index, Owner: from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	// 거스름돈
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{Owner: from, Amount: change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{Owner: to, Amount: amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{ID: "", Timestamp: int(time.Now().Unix()), TxIns: txIns, TxOuts: txOuts}
	return tx, nil
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{Owner: "COINBASE", TxID: "", Index: -1},
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
	tx.getID()
	return &tx
}
