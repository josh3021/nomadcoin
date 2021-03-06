package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"os"

	"github.com/josh3021/nomadcoin/utils"
)

type fileLayer interface {
	hasWalletFile() bool
	writeFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type layer struct{}

func (layer) hasWalletFile() bool {
	_, err := os.Stat(walletFilename)
	return !os.IsNotExist(err)
}

func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

var files fileLayer = layer{}

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

const walletFilename string = "nomadcoin.wallet"

var w *wallet

func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privateKey
}

func persistPrivateKey(privateKeyBytes []byte) {
	err := files.writeFile(walletFilename, privateKeyBytes, 0644)
	utils.HandleErr(err)
}

func readWalletFile() []byte {
	walletBytes, err := files.readFile(walletFilename)
	utils.HandleErr(err)
	return walletBytes
}

func marshalWalletBytes(privateKey *ecdsa.PrivateKey) []byte {
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	utils.HandleErr(err)
	return privateKeyBytes
}

func parseWalletBytes(bytes []byte) *ecdsa.PrivateKey {
	privateKey, err := x509.ParseECPrivateKey(bytes)
	utils.HandleErr(err)
	return privateKey
}

func parseAddress(privateKey *ecdsa.PrivateKey) string {
	return utils.EncodeBigInts(privateKey.X, privateKey.Y)
}

func Sign(payload string, wallet *wallet) string {
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, wallet.privateKey, payloadBytes)
	utils.HandleErr(err)
	// signature
	return utils.EncodeBigInts(r, s)
}

// Verify verfies signature
func Verify(signature, payload, address string) bool {
	r, s, err := utils.RestoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := utils.RestoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	return ecdsa.Verify(&publicKey, payloadBytes, r, s)
}

// Wallet returns wallet (Initialize wallet if it does not initialized).
func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if files.hasWalletFile() {
			walletBytes := readWalletFile()
			privateKey := parseWalletBytes(walletBytes)
			w.privateKey = privateKey
		} else {
			privateKey := createPrivateKey()
			privateKeyBytes := marshalWalletBytes(privateKey)
			persistPrivateKey(privateKeyBytes)
			w.privateKey = privateKey
		}
		w.Address = parseAddress(w.privateKey)
	}
	return w
}
