package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Tx
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(txs []*Tx, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  txs,
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

func NewGenesisBlock(coinbase *Tx) *Block {
	return NewBlock([]*Tx{coinbase}, nil)
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
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
