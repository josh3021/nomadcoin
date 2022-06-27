package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	hash := "ada0db1bcc3bb7d19b34d0911e00ffde1acb4bb632f2f2d6b8b8ced3470ec673"
	s := struct{ Test string }{Test: "text"}
	t.Run("Hash should always be same", func(t *testing.T) {
		r := Hash(s)
		if hash != r {
			t.Errorf("Expected %s, got %s", hash, r)
		}
	})
	t.Run("Hash should be hex encoded", func(t *testing.T) {
		r := Hash(s)
		_, err := hex.DecodeString(r)
		if err != nil {
			t.Error("Hash should be hex encoded")
		}
	})
}

func ExampleHash() {
	s := struct{ Test string }{Test: "text"}
	x := Hash(s)
	fmt.Println(x)
	// Output: ada0db1bcc3bb7d19b34d0911e00ffde1acb4bb632f2f2d6b8b8ced3470ec673
}

func TestToBytes(t *testing.T) {
	s := "test"
	r := ToBytes(s)
	// t.Log(r)
	k := reflect.TypeOf(r).Kind()
	if k != reflect.Slice {
		t.Errorf("ToBytes should return a slice of bytes, but got %s", k)
	}
}

func ExampleToBytes() {
	s := "test"
	r := ToBytes(s)
	fmt.Println(r)
	// Output: [7 12 0 4 116 101 115 116]
}

func TestSplitter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}
	tests := []test{
		{input: "0:6:0", sep: ":", index: 1, output: "6"},
		{input: "0:6:0", sep: ":", index: 10, output: ""},
		{input: "0:6:0", sep: "/", index: 0, output: "0:6:0"},
		{input: "0:6:0", sep: "/", index: 1, output: ""},
	}
	for _, tc := range tests {
		got := Splitter(tc.input, tc.sep, tc.index)
		if got != tc.output {
			t.Errorf("Expected %s, got %s", tc.output, got)
		}
	}
}

func ExampleSplitter() {
	r := Splitter("0:6:0", ":", 1)
	fmt.Println(r)
	// Output: 6
}

func TestHandleErr(t *testing.T) {
	text := "ERROR FOR TEST"
	defer func() {
		if r := recover(); r == nil {
			fmt.Println(r)
			t.Errorf("Expected Panic with %s, Got nil", text)
		}
	}()
	err := errors.New(text)
	HandleErr(err)
}

func TestToJSON(t *testing.T) {
	type testStruct struct{ Text string }
	input := testStruct{"test"}
	r := ToJSON(input)
	t.Logf("%s", r)
	kind := reflect.TypeOf(r).Kind()
	if kind != reflect.Slice {
		t.Errorf("Expected %v, got %v", reflect.Slice, kind)
	}
	var restored testStruct
	json.Unmarshal(r, &restored)
	if !reflect.DeepEqual(restored, input) {
		t.Errorf("ToJSON() should encode correctly")
	}
}

func ExampleToJSON() {
	type testStruct struct{ Text string }
	input := testStruct{"test"}
	r := ToJSON(input)
	fmt.Printf("%s", r)
	// Output: {"Text":"test"}
}

func TestFromBytes(t *testing.T) {
	type testStruct struct{ Text string }
	input := testStruct{"test"}
	j := ToBytes(input)
	var restored testStruct
	FromBytes(&restored, j)
	if !reflect.DeepEqual(input, restored) {
		t.Error("FromBytes() should restore the struct.")
	}
}

func ExampleFromBytes() {
	type testStruct struct{ Text string }
	var restored testStruct
	i := testStruct{Text: "test"}
	j := ToBytes(i)
	FromBytes(&restored, j)
	fmt.Println(restored)
	// Output: {test}
}

func TestEncodeBigInts(t *testing.T) {
	bigA := big.NewInt(123123)
	bigB := big.NewInt(456456)
	expected := "01e0f306f708"
	s := EncodeBigInts(bigA, bigB)
	if s != expected {
		t.Errorf("Expected: %s, got: %s", expected, s)
	}
}

func ExampleEncodeBigInts() {
	bigA := big.NewInt(123123)
	bigB := big.NewInt(456456)
	s := EncodeBigInts(bigA, bigB)
	fmt.Println(s)
	// Output: 01e0f306f708
}

func TestRestoreBigInts(t *testing.T) {
	t.Run("RestoreBigInts should return two-seperated bigInts", func(t *testing.T) {
		i := "01e0f306f708"
		var expectedBigA int64 = 123123
		var expectedBigB int64 = 456456
		bigA, bigB, _ := RestoreBigInts(i)
		if !reflect.DeepEqual(bigA, big.NewInt(expectedBigA)) {
			t.Errorf("Expected: %v, got: %v", expectedBigA, bigA)
		}
		if !reflect.DeepEqual(bigB, big.NewInt(expectedBigB)) {
			t.Errorf("Expected: %v, got: %v", expectedBigB, bigB)
		}
	})
	t.Run("RestoreBigInts should return ERROR when Input is invalid", func(t *testing.T) {
		i := "0245l2"
		_, _, err := RestoreBigInts(i)
		if err == nil {
			t.Error("Expected ERROR, got nil")
		}
	})
}

func ExampleRestoreBigInts() {
	i := "01e0f306f708"
	bigA, bigB, _ := RestoreBigInts(i)
	fmt.Println(bigA, bigB)
	// Output: 123123 456456
}
