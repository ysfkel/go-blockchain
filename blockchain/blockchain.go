package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath        = "./tmp/blocks"
	latestHashKey = "lh"
)

type IBlockChain interface {
	AddBlock(data string)
	GetLastestBlock() *Block
}

type BlockChain struct {
	LastHash []byte // hash of the last block
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func createBlock(data string, prevHash []byte) *Block {
	initialNonce := 0
	block := &Block{[]byte{}, []byte(data), prevHash, initialNonce}

	//create proof
	pow := NewProof(block)
	//run proof of work
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	var latestHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(latestHashKey))
		HandleError(err)
		latestHash, err = item.ValueCopy(nil)

		return err
	})

	HandleError(err)

	newBlock := createBlock(data, latestHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(err)
		err = txn.Set([]byte(latestHashKey), newBlock.Hash)
		return err
	})

	HandleError(err)
}

func Genesis() *Block {
	return createBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {

	var lastHash []byte
	//open db
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		// check if we have a d
		_, err := txn.Get([]byte(latestHashKey))

		if err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			//generate genesis block
			genesis := Genesis()
			fmt.Println("Genesis proved")
			//add genesis to database
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleError(err)
			//set genesis as last hash
			err = txn.Set([]byte(latestHashKey), genesis.Hash)

			lastHash = genesis.Hash

			return err

		}

		item, err := txn.Get([]byte(latestHashKey))
		HandleError(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})

	HandleError(err)

	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

func (chain *BlockChain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}

	return iter
}

//iterates backwards from the newest to genesis
func (iter *BlockchainIterator) Next() *Block {

	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		HandleError(err)
		encodedBlock, err := item.ValueCopy(nil)

		block = BlockBytes(encodedBlock).Deserialize()

		return err
	})

	HandleError(err)
	// change iterator hash to previous hash
	iter.CurrentHash = block.PrevHash

	return block
}
