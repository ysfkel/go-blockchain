package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
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

	HandleError(err)

	return res.Bytes()
}

func (b *Block) SerializeToBlockBytes() BlockBytes {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	HandleError(err)

	return BlockBytes(res.Bytes())
}

func (b BlockBytes) Deserialize() *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(b))

	err := decoder.Decode(&block)

	HandleError(err)

	return &block
}

func (b BlockBytes) ToBytes() []byte {
	return []byte(b)
}

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
