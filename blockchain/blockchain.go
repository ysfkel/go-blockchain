package block 

type IBlockChain interface {
	  AddBlock(data string)
	  GetLastestBlock() *Block  
}

type BlockChain struct {
	Blocks []*Block
}
 
func createBlock(data string, prevHash []byte) * Block {
	initialNonce := 0
	block := &Block{[]byte{}, []byte(data), prevHash, initialNonce}
	
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (chain *BlockChain) GetLastestBlock() *Block {
	if len(chain.Blocks) == 0 {
		return nil
	}

	return chain.Blocks[len(chain.Blocks) - 1]

}

func (chain *BlockChain) AddBlock(data string) {
     prevBlock := chain.GetLastestBlock()
     newBlock := createBlock(data, prevBlock.Hash)

	 chain.Blocks = append(chain.Blocks, newBlock)
}

func Genesis() * Block {
	return createBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

 