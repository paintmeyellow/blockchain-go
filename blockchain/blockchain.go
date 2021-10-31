package blockchain

import "github.com/boltdb/bolt"

type Blockchain struct {
	tip []byte
	DB  *bolt.DB
}

const (
	dbFile       = "bolt-blockchain"
	blocksBucket = "blocks"
)

func NewBlockchain() *Blockchain {
	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)
	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b != nil {
			tip = b.Get([]byte("l"))
			return nil
		}
		genesis := NewGenesisBlock()
		b, _ = tx.CreateBucket([]byte(blocksBucket))
		_ = b.Put(genesis.Hash, genesis.Serialize())
		_ = b.Put([]byte("l"), genesis.Hash)
		tip = genesis.Hash
		return nil
	})
	return &Blockchain{tip: tip, DB: db}
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	_ = bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	newBlock := NewBlock(data, lastHash)
	_ = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		_ = b.Put(newBlock.Hash, newBlock.Serialize())
		_ = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})
}

func (bc *Blockchain) Iterator() *Iterator {
	return &Iterator{
		currentHash: bc.tip,
		db:          bc.DB,
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
