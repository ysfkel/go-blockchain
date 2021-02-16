package block 

import (
	"math"
	"fmt"
	"math/big"
	"bytes"
	"encoding/binary"
	"log"
	"crypto/sha256"
)

/*
 1. Get the data from the block
 2. Create a counter (nonce) which starts at 0
 3. Create a hash of the data plus the counter
 4. Check the resulting hash to see if it meets a set of requirements 
    (this is the idae of difficulty)
 5. if the hash meets the set of requirements then we use that hash and say it signs the block
 6. otherwise we go back and create another hash and we repeat the process until we get a 
    hash that meets the set of requirements
  
Requirement 
First few bytes of the hash must contain 0s
*/

/*
The difficulty is a measure of how difficult it is to mine a Bitcoin block, or in more technical terms, to find a hash below a given target.
*/
const Difficulty uint = 18

type ProofOfWork struct {
	Block *Block 
	//a target hash is a numeric value that a hashed block header must be less than or equal to in order for a new block to be awarded to a miner. 
	Target *big.Int // represents the requirement described which is derived from the difficuly
}

/*
 Takes a pointer to a block 
 And then produce a pointer to a proof of work
 - takes the block and pairs with the target
*/
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	/*
	 substract difficulty from 256 (number of bytes in the hashes)
	 left shift the number of bytes by the target
	*/
	result := uint(256-Difficulty)
	target.Lsh(target, result)

	pow := &ProofOfWork{b, target}

	return pow
}

/*

**/
func (pow *ProofOfWork) InitData(nonce int) []byte{

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

//4. Check the resulting hash to see if it meets a set of requirements 
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int 
	var hash [32]byte 

	nonce := 0 

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		//hash data
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		//convert hash to big interger
		intHash.SetBytes(hash[:])
		//compare the big int version of hash to target
        if intHash.Cmp(pow.Target) == -1 {
		   //means hash is less than target , meaning we have 
		   //sgned the block
           break
		} else {
           nonce++
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int 

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}




