package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/josh3021/nomadcoin/utils"
	"github.com/josh3021/nomadcoin/wallet"
)

const minerReward int = 50

// Tx contains information of transactions
type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

// TxIn contains information of transactions input
type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

// TxOut contains information of transactions Output
type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

// UTxOut contains information of Unconfirmed transactions Output
type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (tx *Tx) getID() {
	tx.ID = utils.Hash(tx)
}

func (tx *Tx) sign() {
	for _, txIn := range tx.TxIns {
		txIn.Signature = wallet.Sign(tx.ID, wallet.Wallet())
	}
}

type mempool struct {
	// Txs []*Tx `json:"txs"`
	Txs map[string]*Tx `json:"txs"`
	m   sync.Mutex
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	// m.Txs = append(m.Txs, tx)
	m.Txs[tx.ID] = tx
	return tx, nil
}

func (m *mempool) ConfirmTxs() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx)
	return txs
}

// Mempool contains not confirmed transactions
var m *mempool
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{Txs: make(map[string]*Tx)}
	})
	return m
}

func MempoolStatus() *mempool {
	mem := Mempool()
	mem.m.Lock()
	defer mem.m.Unlock()

	return mem
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool().Txs {
		for _, txIn := range tx.TxIns {
			if txIn.TxID == uTxOut.TxID && txIn.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			break
		}
	}
	return valid
}

var errorNotEnoghMoney = errors.New("not enough money")
var errorTxNotValid = errors.New("tx not valid")

func makeTx(from, to string, amount int) (*Tx, error) {
	if GetBalanceByAddress(Blockchain(), from) < amount {
		return nil, errorNotEnoghMoney
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	uTxOuts := UTxOutsByAddress(Blockchain(), from)
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{TxID: uTxOut.TxID, Index: uTxOut.Index, Signature: from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	// 거스름돈
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{Address: from, Amount: change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{Address: to, Amount: amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{ID: "", Timestamp: int(time.Now().Unix()), TxIns: txIns, TxOuts: txOuts}
	tx.getID()
	tx.sign()
	if !validate(tx) {
		return nil, errorTxNotValid
	}
	return tx, nil
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{Signature: "COINBASE", TxID: "", Index: -1},
	}
	txOuts := []*TxOut{
		{Address: address, Amount: minerReward},
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

func (b *blockchain) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()
	// m.Txs = append(m.Txs, tx)
	m.Txs[tx.ID] = tx
}
