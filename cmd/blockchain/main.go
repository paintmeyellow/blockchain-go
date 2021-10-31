package main

import (
	"demo/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.DB.Close()

	cli := blockchain.CLI{BC: bc}
	cli.Run()
}
