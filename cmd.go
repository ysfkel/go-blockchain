package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/ysfkel/go-blockchain/blockchain"
	blc "github.com/ysfkel/go-blockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("add -block BLOCK DATA - add a block to the chain")
	fmt.Println("print - Prints the blocks in the chain")
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) printchain() {

	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("%s\n", block.Data)
		fmt.Printf("%x\n", block.Hash)
		fmt.Printf("%x\n", block.PrevHash)
		fmt.Printf("-----------Validate block-----------\n")
		pow := blc.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validdateArgs()
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse((os.Args[2:]))
		blockchain.HandleError(err)
	case "print":
		err := printChainCmd.Parse((os.Args[2:]))
		blockchain.HandleError(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printUsage()
	}

}
func (cli *CommandLine) validdateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		// exits application by shutting down the go routine
		// initaites a shutdow
		runtime.Goexit()
	}
}
