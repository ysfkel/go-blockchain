package main

import (
	"os"

	blc "github.com/ysfkel/go-blockchain/blockchain"
)

func main() {
	defer os.Exit(0)
	chain := blc.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()
}
