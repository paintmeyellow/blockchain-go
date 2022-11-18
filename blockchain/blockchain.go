package blockchain

import (
	"encoding/hex"

	"github.com/boltdb/bolt"
)

const (
	blocksBucket        = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func Create(addr string, db *bolt.DB) (*Blockchain, error) {
	var tip []byte
	err := db.Update(func(tx *bolt.Tx) error {
		coinbase, err := NewCoinbaseTx(addr, genesisCoinbaseData)
		if err != nil {
			return err
		}
		genesis := NewGenesisBlock(coinbase)
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}
		if err = b.Put(genesis.Hash, genesis.Serialize()); err != nil {
			return err
		}
		if err = b.Put([]byte("l"), genesis.Hash); err != nil {
			return err
		}
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Blockchain{tip: tip, db: db}, nil
}

func Open(db *bolt.DB) (*Blockchain, error) {
	var tip []byte
	err := db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(blocksBucket)); b != nil {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Blockchain{tip: tip, db: db}, nil
}

func (bc *Blockchain) MineBlock(txs []*Tx) error {
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

func (bc *Blockchain) UnspentTxs(address string) []*Tx {
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

func (bc *Blockchain) UTXO(address string) []TxOutput {
	var outs []TxOutput
	txs := bc.UnspentTxs(address)
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
func (bc *Blockchain) SpendableOutputs(addr string, amount int) (acc int, utxo map[string][]int) {
	utxo = make(map[string][]int)
	unspentTXs := bc.UnspentTxs(addr)
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

func (bc *Blockchain) Close() {
	if bc.db != nil {
		if err := bc.db.Close(); err != nil {

		}
	}
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
