package usecase

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"blockchain-go/blockchain"
)

var errNotEnoughFunds = errors.New("blockchain: not enough funds")

type PayToBC interface {
	SpendableOutputs(ctx context.Context, addr string, amount int) (acc int, utxo map[string][]int)
	MineBlock(ctx context.Context, txs []*blockchain.Tx) error
}

type PayToUcase struct {
	bc PayToBC
	tr trace.Tracer
}

func NewPayToUcase(bc PayToBC) *PayToUcase {
	return &PayToUcase{
		bc: bc,
		tr: otel.Tracer("usecase"),
	}
}

func (ucase PayToUcase) Handle(ctx context.Context, from, to string, amount int) error {
	ctx, span := ucase.tr.Start(ctx, "usecase.pay_to")
	defer span.End()

	acc, utxo := ucase.bc.SpendableOutputs(ctx, from, amount)
	if acc < amount {
		span.RecordError(errNotEnoughFunds)
		return errNotEnoughFunds
	}
	tx, err := blockchain.NewTx(from, to, amount, acc, utxo)
	if err != nil {
		span.RecordError(err)
		return err
	}
	span.AddEvent("mining block")
	if err = ucase.bc.MineBlock(ctx, []*blockchain.Tx{tx}); err != nil {
		span.RecordError(err)
		return err
	}
	span.AddEvent("block has been mined")
	return nil
}
