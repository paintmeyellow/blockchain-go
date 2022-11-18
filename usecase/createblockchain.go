package usecase

import (
	"demo/blockchain"

	"github.com/boltdb/bolt"
)

type CreateBlockchainUcase struct {
	db *bolt.DB
}

func NewCreateBlockchainUcase(db *bolt.DB) *CreateBlockchainUcase {
	return &CreateBlockchainUcase{db: db}
}

func (ucase *CreateBlockchainUcase) Handle(addr string) error {
	_, err := blockchain.Create(addr, ucase.db)
	return err
}
