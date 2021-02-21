package main

import (
	"os"

	blc "github.com/ysfkel/go-blockchain/blockchain"
	db "github.com/ysfkel/go-blockchain/blockchain/db"
	cmd "github.com/ysfkel/go-blockchain/cmd"
)

func main() {
	defer os.Exit(0)
	chain := blc.InitBlockChain(db.NewRepository())
	// close db after go channel exits properly (runtime.Goexit())
	defer chain.Database.Close()

	cli := cmd.CommandLine{chain}
	cli.Run()
}
