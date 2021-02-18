package block

import (
	"bytes"
	"encoding/gob"

	shared "github.com/ysfkel/go-blockchain/shared"
)

type BlockBytes []byte

type Block struct {
	Hash     []byte //current block hash: hash(Data, PrevHash)
	Data     []byte
	PrevHash []byte //last block hash
	Nonce    int
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
