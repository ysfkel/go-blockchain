package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
	blk "github.com/ysfkel/go-blockchain/blockchain/block"
	"github.com/ysfkel/go-blockchain/blockchain/consensus"
	"github.com/ysfkel/go-blockchain/blockchain/db"
	"github.com/ysfkel/go-blockchain/shared"
)

type IBlockChain interface {
	AddBlock(data string)
	GetLastestBlock() *blk.Block
}

type BlockChain struct {
	LastHash []byte // hash of the last block
	Database db.IRepository
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    db.IRepository
}

func createBlock(data string, prevHash []byte) *blk.Block {
	initialNonce := 0
	block := &blk.Block{[]byte{}, []byte(data), prevHash, initialNonce}

	//create proof
	pow := consensus.NewProof(block)
	//run proof of work
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (chain *BlockChain) AddBlock(data string) {

	latestHash, err := chain.Database.Get([]byte(shared.LatestHashKey))

	shared.HandleError(err)

	newBlock := createBlock(data, latestHash)

	err = chain.Database.Update(newBlock)

	shared.HandleError(err)
}

func Genesis() *blk.Block {
	return createBlock(shared.Genesis, []byte{})
}

func InitBlockChain() *BlockChain {

	repo := db.NewRepository()

	lastHash, err := repo.Get([]byte(shared.LatestHashKey))

	if err == badger.ErrKeyNotFound {
		fmt.Println("No existing blockchain found")
		fmt.Println("Generating genesis block")
		genesis := Genesis()
		err = repo.Update(genesis)
		shared.HandleError(err)
		lastHash = genesis.Hash
	}

	blockchain := BlockChain{
		LastHash: lastHash,
		Database: repo,
	}

	return &blockchain
}

func (chain *BlockChain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}

	return iter
}

//iterates backwards from the newest to genesis
func (iter *BlockchainIterator) Next() *blk.Block {

	var block *blk.Block

	encodedBlock, err := iter.Database.Get(iter.CurrentHash)

	shared.HandleError(err)

	block = blk.BlockBytes(encodedBlock).Deserialize()
	// change iterator hash to previous hash
	iter.CurrentHash = block.PrevHash

	return block
}
