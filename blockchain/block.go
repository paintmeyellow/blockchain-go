package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", nil)
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	_ = gob.NewEncoder(&res).Encode(b)
	return res.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	_ = gob.NewDecoder(bytes.NewReader(d)).Decode(&block)
	return &block
}
