package usecase

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"blockchain-go/blockchain"
)

type GetBalanceBC interface {
	UTXO(ctx context.Context, addr string) []blockchain.TxOutput
}

type GetBalanceUcase struct {
	bc GetBalanceBC
	tr trace.Tracer
}

func NewGetBalanceUcase(bc GetBalanceBC) *GetBalanceUcase {
	return &GetBalanceUcase{
		bc: bc,
		tr: otel.Tracer("usecase"),
	}
}

type Balance struct {
	Value int
}

func (ucase *GetBalanceUcase) Handle(ctx context.Context, addr string) *Balance {
	ctx, span := ucase.tr.Start(ctx, "GetBalanceUcase.Handle")
	defer span.End()

	var balance int
	utxo := ucase.bc.UTXO(ctx, addr)
	for _, out := range utxo {
		balance += out.Value
	}
	return &Balance{Value: balance}
}
