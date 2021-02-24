package transaction

import (
	"bytes"
	"fmt"
)

type Transaction struct {
	ID       []byte
	Inputs   []TxInput
	TxOutput []TxOutput
}

type TxOutput struct {
	Value  int    // value that is locked in the output / an output is indivisible
	PubKey string // needed to unlock that tokens in the value field,
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string // provdies data which is used to the outputs pub key
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	//encode :=
}

func CoinbaseTx(to, data string) *Transaction {

	if data == "" {
		data = fmt.Sprintf("cOINS TO %S", to)
	}

	txin := TxInput{
		ID:  []byte{}, //coinbase tx references no output
		Out: -1,       // -1 because coinbase tx references no output
		Sig: data,
	}

	txout := TxOutput{
		Value:  100,
		PubKey: to, // if an account named jack mines this block, then jack gets these tokens
	}

	tx := Transaction{
		nil,
		[]TxInput{txin},
		[]TxOutput{txout},
	}
	tx.SetID()

	return &tx

}
