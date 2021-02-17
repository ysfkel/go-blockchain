package blockchain

const (
	dbPath        = "./tmp/blocks"
	latestHashKey = "lh"
	genesis       = "Genesis"
)

type IBlockChain interface {
	AddBlock(data string)
	GetLastestBlock() *Block
}

type BlockChain struct {
	LastHash []byte      // hash of the last block
	Database IRepository //*badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    IRepository
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

	latestHash, err := chain.Database.Get([]byte(latestHashKey))

	HandleError(err)

	newBlock := createBlock(data, latestHash)

	err = chain.Database.Update(newBlock)

	HandleError(err)
}

func Genesis() *Block {
	return createBlock(genesis, []byte{})
}
func InitBlockChain() *BlockChain {

	repo, lastHash := initRepository()

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
func (iter *BlockchainIterator) Next() *Block {

	var block *Block

	encodedBlock, err := iter.Database.Get(iter.CurrentHash)

	HandleError(err)

	block = BlockBytes(encodedBlock).Deserialize()
	// change iterator hash to previous hash
	iter.CurrentHash = block.PrevHash

	return block
}
