package blockchain

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Tx
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	tr            trace.Tracer
}

func NewBlock(txs []*Tx, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  txs,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		tr:            otel.Tracer("blockchain"),
	}
	return block
}

func NewGenesisBlock(coinbase *Tx) *Block {
	b := NewBlock([]*Tx{coinbase}, nil)
	return b
}

func (b *Block) Mine(ctx context.Context) {
	ctx, span := b.tr.Start(ctx, "Block.Mine")
	defer span.End()
	pow := NewProofOfWork(b)
	nonce, hash := pow.Run(ctx)
	b.Hash = hash[:]
	b.Nonce = nonce
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
