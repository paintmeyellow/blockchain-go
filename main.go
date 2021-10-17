package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

type Blockchain struct {
	blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) AddBlock(b *Block) {
	if b == nil {
		return
	}
	prevBlock := bc.blocks[len(bc.blocks)-1]
	block := *b
	block.PrevBlockHash = prevBlock.Hash
	bc.blocks = append(bc.blocks, &block)
}

func NewBlock(data []byte) *Block {
	b := Block{
		Timestamp:     time.Now().Unix(),
		Data:          data,
		PrevBlockHash: nil,
		Hash:          nil,
	}
	b.SetHash()
	return &b
}

func NewGenesisBlock() *Block {
	return NewBlock([]byte("Genesis Block"))
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func main() {
	bc := NewBlockchain()
	bc.AddBlock(NewBlock([]byte("Send 1 BTC to Ivan")))
	bc.AddBlock(NewBlock([]byte("Send 2 BTC to Ivan")))

	for _, block := range bc.blocks {
		fmt.Printf("Prev: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Println()
	}
}
