package main 

import (
	"fmt"
	"strconv"
	blc "github.com/ysfkel/go-blockchain/blockchain"
)
 
func main() {
   chain := blc.InitBlockChain()
   chain.AddBlock("First Block after Genesis")
   chain.AddBlock("Seconds Block after Genesis")
   chain.AddBlock("Third Block after Genesis")


   for _, block := range chain.Blocks {
	  fmt.Printf("%s\n", block.Data)
	  fmt.Printf("%x\n", block.Hash)
	  fmt.Printf("%x\n", block.PrevHash)
	  fmt.Printf("-----------Validate block-----------\n")
	  pow := blc.NewProof(block)
	  fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	  fmt.Println()
   }
}