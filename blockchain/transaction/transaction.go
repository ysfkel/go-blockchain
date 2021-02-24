package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/ysfkel/go-blockchain/shared"
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

	encode := gob.NewEncoder(&encoded)

	err := encode.Encode(tx)
	shared.HandleError(err)

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]
}

func CoinbaseTx(to, data string) *Transaction {

	if data == "" {
		data = fmt.Sprintf("Coins TO %s", to)
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

func (tx *Transaction) IsCoinbase() bool {

	return len(tx.Inputs) == 1 &&
		len(tx.Inputs[0].ID) == 0 &&
		tx.Inputs[0].Out == -1
}

//checks if the account (data) can unlock
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

//checks if the account (data) owns the information
//or owns the information in the output which is referenced
//by the input
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
