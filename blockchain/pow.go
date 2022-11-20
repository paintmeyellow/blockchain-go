package blockchain

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const targetBits = 20

type ProofOfWork struct {
	block  *Block
	target *big.Int
	tr     trace.Tracer
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{
		block:  b,
		target: target,
		tr:     otel.Tracer("pow"),
	}
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.HashTransactions(),
		[]byte(strconv.Itoa(int(pow.block.Timestamp))),
		[]byte(strconv.Itoa(targetBits)),
		[]byte(strconv.Itoa(nonce)),
	}, []byte{})
}

func (pow *ProofOfWork) Run(ctx context.Context) (int, []byte) {
	_, span := pow.tr.Start(ctx, "ProofOfWork.Run")
	defer span.End()

	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}
