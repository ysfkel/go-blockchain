package blockchain

import (
	"fmt"
	"runtime"

	"github.com/dgraph-io/badger"
	blk "github.com/ysfkel/go-blockchain/blockchain/block"
	"github.com/ysfkel/go-blockchain/blockchain/consensus"
	"github.com/ysfkel/go-blockchain/blockchain/db"
	tx "github.com/ysfkel/go-blockchain/blockchain/transaction"
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

func createBlock(transactions []*tx.Transaction, prevHash []byte) *blk.Block {
	initialNonce := 0
	block := &blk.Block{[]byte{}, transactions, prevHash, initialNonce}

	//create proof
	pow := consensus.NewProof(block)
	//run proof of work
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (chain *BlockChain) AddBlock(transactions []*tx.Transaction) {

	latestHash, err := chain.Database.Get([]byte(shared.LatestHashKey))

	shared.HandleError(err)

	newBlock := createBlock(transactions, latestHash)

	err = chain.Database.Update(newBlock)

	shared.HandleError(err)
}

func createGenesis(coinbase *tx.Transaction) *blk.Block {
	return createBlock([]*tx.Transaction{coinbase}, []byte{})
}

func InitBlockChain(address string, repo db.IRepository) *BlockChain {

	if repo.DBExists() == false {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	lastHash, err := repo.Get([]byte(shared.LatestHashKey))

	if err == badger.ErrKeyNotFound {
		fmt.Println("No existing block found")
		fmt.Println("Generating genesis block")
		genesis := createGenesis(tx.CoinbaseTx("yusuf-address", "coinbase tx"))
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

func (chain *BlockChain) GetBlockChain(address string, repo db.IRepository) *BlockChain {

	if repo.DBExists() == false {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	lastHash, err := repo.Get([]byte(shared.LatestHashKey))

	shared.HandleError(err)

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
