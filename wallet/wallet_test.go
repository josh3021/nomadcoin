package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"

	"github.com/josh3021/nomadcoin/utils"
)

const (
	testKey       string = "307702010104209f2a3398b92ba15fd6224f77afab9b3cf26d90c3f0cb0a86f683a6d1c0a033d0a00a06082a8648ce3d030107a14403420004495566a10c82156b2ac498395ae3c31634add7351272ca6d20a749ec7ae8e46ff5e877d9030ee1b13771badbc33b6ad1a36efe4838a207b3fd87d441a3f72d65"
	testPayload   string = "765cd5fbc8bdb14616f299bb4e7952065cd5d153d3511e680d5ddd7cc74e1aaf"
	testSignature string = "b4a402235d90ecbf90c4c41c13614261dc4c0b311bbe9499537f67b9a22cfbae6edeb0da3e58cb77596f02623ab2eada22d272cc347c9514a029c97631a449f1"
)

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

func makeTestWallet() *wallet {
	w := &wallet{}
	b, err := hex.DecodeString(testKey)
	utils.HandleErr(err)
	pk, err := x509.ParseECPrivateKey(b)
	utils.HandleErr(err)
	w.privateKey = pk
	w.Address = parseAddress(w.privateKey)
	return w
}

func TestSign(t *testing.T) {
	w := makeTestWallet()
	s := Sign(testPayload, w)
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign should return hex-encoded string, got: %s", s)
	}
}

func TestVerify(t *testing.T) {
	w := makeTestWallet()
	t.Run("Verify should have correct payload.", func(t *testing.T) {
		v := Verify(testSignature, testPayload, w.Address)
		if !v {
			t.Error("Verify should have correct payload.")
		}
	})
	t.Run("Verify should return error if signature is invalid.", func(t *testing.T) {
		v := Verify("b4a402235d90ecbf90c4c41c13614261dc4c0b311bbe9499537f67b9a22cfbae6edeb0da3e58cb77596f02623ab2eada22d272cc347c9514a029c97631a449f2", testPayload, w.Address)
		if v {
			t.Error("Expected: false, got: true")
		}
	})
	t.Run("Verify should return error if payload is invalid", func(t *testing.T) {
		v := Verify(testSignature, "765cd5fbc8bdb14616f299bb4e7952065cd5d153d3511e680d5ddd7cc74e1aae", w.Address)
		if v {
			t.Error("Expected: false, got: true")
		}
	})
}

func TestWallet(t *testing.T) {
	t.Run("Create New Wallet", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool {
				return false
			},
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})
	t.Run("Restore Existing Wallet", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool {
				return true
			},
		}
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})
}

// func TestParseWalletBytes(t *testing.T) {
// 	b, err := hex.DecodeString(testKey)
// 	utils.HandleErr(err)
// 	pk := parseWalletBytes(b)
// }
