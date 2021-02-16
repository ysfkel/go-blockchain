package block 


type Block struct {
	Hash []byte //current block hash: hash(Data, PrevHash)
	Data []byte 
	PrevHash []byte //last block hash
	Nonce int
}
  

 