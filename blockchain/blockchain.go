package blockchain

import (
	"context"
	"encoding/hex"

	"github.com/boltdb/bolt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	blocksBucket        = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
	tr  trace.Tracer
}

func New(db *bolt.DB) *Blockchain {
	return &Blockchain{
		db: db,
		tr: otel.Tracer("blockchain"),
	}
}

func (bc *Blockchain) Create(ctx context.Context, addr string) error {
	ctx, span := bc.tr.Start(ctx, "Blockchain.Create")
	defer span.End()

	var tip []byte
	err := bc.db.Update(func(tx *bolt.Tx) error {
		coinbase, err := NewCoinbaseTx(addr, genesisCoinbaseData)
		if err != nil {
			span.SetStatus(codes.Error, "new coinbase tx")
			span.RecordError(err)
			return err
		}
		genesis := NewGenesisBlock(coinbase)
		genesis.Mine(ctx)
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			span.SetStatus(codes.Error, "create bucket")
			span.RecordError(err)
			return err
		}
		if err = b.Put(genesis.Hash, genesis.Serialize()); err != nil {
			span.SetStatus(codes.Error, "write block data")
			span.RecordError(err)
			return err
		}
		if err = b.Put([]byte("l"), genesis.Hash); err != nil {
			span.SetStatus(codes.Error, "write block header")
			span.RecordError(err)
			return err
		}
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		span.SetStatus(codes.Error, "execute db transaction")
		span.RecordError(err)
		return err
	}
	bc.tip = tip
	return nil
}

func (bc *Blockchain) Open() error {
	var tip []byte
	err := bc.db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(blocksBucket)); b != nil {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	bc.tip = tip
	return nil
}

func (bc *Blockchain) MineBlock(ctx context.Context, txs []*Tx) error {
	ctx, span := bc.tr.Start(ctx, "Blockchain.MineBlock")
	defer span.End()
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		return err
	}
	newBlock := NewBlock(txs, lastHash)
	newBlock.Mine(ctx)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if err = b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			return err
		}
		if err = b.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}
		bc.tip = newBlock.Hash
		return nil
	})
	return err
}

func (bc *Blockchain) UnspentTxs(ctx context.Context, address string) []*Tx {
	_, span := bc.tr.Start(ctx, "Blockchain.UnspentTxs")
	defer span.End()
	var unspentTXs []*Tx
	var block *Block
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for {
		block = bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Vout {
				for _, spentOut := range spentTXOs[txID] {
					if spentOut == outIdx {
						continue Outputs
					}
				}
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

func (bc *Blockchain) UTXO(ctx context.Context, address string) []TxOutput {
	ctx, span := bc.tr.Start(ctx, "Blockchain.UTXO")
	defer span.End()
	var outs []TxOutput
	txs := bc.UnspentTxs(ctx, address)
	for _, tx := range txs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				outs = append(outs, out)
			}
		}
	}
	return outs
}

// SpendableOutputs returns map[txID][]vout
func (bc *Blockchain) SpendableOutputs(ctx context.Context, addr string, amount int) (acc int, utxo map[string][]int) {
	ctx, span := bc.tr.Start(ctx, "Blockchain.SpendableOutputs")
	defer span.End()
	utxo = make(map[string][]int)
	unspentTXs := bc.UnspentTxs(ctx, addr)
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(addr) && acc < amount {
				acc += out.Value
				utxo[txID] = append(utxo[txID], outIdx)
				if acc >= amount {
					return
				}
			}
		}
	}
	return
}

func (bc *Blockchain) Iterator() *Iterator {
	return &Iterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
}

type Iterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *Iterator) Next() *Block {
	var block *Block
	_ = i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	i.currentHash = block.PrevBlockHash
	return block
}
