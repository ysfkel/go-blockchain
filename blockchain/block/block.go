package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	tx "github.com/ysfkel/go-blockchain/blockchain/transaction"
	shared "github.com/ysfkel/go-blockchain/shared"
)

type BlockBytes []byte

type Block struct {
	Hash         []byte //current block hash: hash(Data, PrevHash)
	Transactions []*tx.Transaction
	PrevHash     []byte //last block hash
	Nonce        int
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	shared.HandleError(err)

	return res.Bytes()
}

func (b *Block) SerializeToBlockBytes() BlockBytes {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	shared.HandleError(err)

	return BlockBytes(res.Bytes())
}

func (b BlockBytes) Deserialize() *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(b))

	err := decoder.Decode(&block)

	shared.HandleError(err)

	return &block
}

func (b BlockBytes) ToBytes() []byte {
	return []byte(b)
}
