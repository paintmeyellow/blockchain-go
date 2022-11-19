package usecase

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"blockchain-go/blockchain"
)

type GetBalanceBC interface {
	UTXO(ctx context.Context, addr string) []blockchain.TxOutput
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
	_, span := otel.Tracer("usecase").Start(ctx, "usecase.get_balance")
	defer span.End()
	span.SetAttributes(attribute.String("addr", addr))

	var balance int
	utxo := ucase.bc.UTXO(ctx, addr)
	for _, out := range utxo {
		balance += out.Value
	}
	return &Balance{Value: balance}
}
