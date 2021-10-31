package main

import (
	"context"
	"demo/blockchain"
	"log"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.DB.Close()

	cli := blockchain.NewCLI(bc)
	if err := cli.Run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
