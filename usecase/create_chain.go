package usecase

import (
	"github.com/boltdb/bolt"

	"blockchain-go/blockchain"
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
