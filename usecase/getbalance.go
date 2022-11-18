package usecase

import (
	"context"
	"demo/blockchain"
)

type GetBalanceBC interface {
	UTXO(addr string) []blockchain.TxOutput
}

type GetBalanceUcase struct {
	bc GetBalanceBC
}

func NewGetBalanceUcase(bc GetBalanceBC) *GetBalanceUcase {
	return &GetBalanceUcase{bc: bc}
}

type Balance struct {
	Value int
}

func (ucase *GetBalanceUcase) Handle(ctx context.Context, addr string) *Balance {
	var balance int
	utxo := ucase.bc.UTXO(addr)
	for _, out := range utxo {
		balance += out.Value
	}
	return &Balance{Value: balance}
}
