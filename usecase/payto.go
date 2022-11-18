package usecase

import (
	"context"
	"errors"

	"blockchain-go/blockchain"
)

var errNotEnoughFunds = errors.New("blockchain: not enough funds")

type PayToBC interface {
	SpendableOutputs(addr string, amount int) (acc int, utxo map[string][]int)
	MineBlock(txs []*blockchain.Tx) error
}

type PayToUcase struct {
	bc PayToBC
}

func NewPayToUcase(bc PayToBC) *PayToUcase {
	return &PayToUcase{bc: bc}
}

func (ucase PayToUcase) Handle(ctx context.Context, from, to string, amount int) error {
	acc, utxo := ucase.bc.SpendableOutputs(from, amount)
	if acc < amount {
		return errNotEnoughFunds
	}
	tx, err := blockchain.NewTx(from, to, amount, acc, utxo)
	if err != nil {
		return err
	}
	if err = ucase.bc.MineBlock([]*blockchain.Tx{tx}); err != nil {
		return err
	}
	return nil
}
