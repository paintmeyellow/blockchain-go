package main

import (
	"context"
	"demo/blockchain"
	"log"
)

func main() {
	cli := blockchain.NewCLI()
	if err := cli.Run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
