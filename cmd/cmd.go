package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/ysfkel/go-blockchain/blockchain"
	"github.com/ysfkel/go-blockchain/blockchain/consensus"
	tx "github.com/ysfkel/go-blockchain/blockchain/transaction"
	"github.com/ysfkel/go-blockchain/shared"
)

type CommandLine struct {
	Blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("add -block BLOCK DATA - add a block to the chain")
	fmt.Println("print - Prints the blocks in the chain")
}

func (cli *CommandLine) addBlock(tx []*tx.Transaction) {
	cli.Blockchain.AddBlock(tx)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) printchain() {

	iter := cli.Blockchain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("-----------Validate block-----------\n")
		pow := consensus.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs()
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse((os.Args[2:]))
		shared.HandleError(err)
	case "print":
		err := printChainCmd.Parse((os.Args[2:]))
		shared.HandleError(err)
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
		cli.printchain()
	}

}
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		// exits application by shutting down the go routine
		// initaites a shutdow
		runtime.Goexit()
	}
}
