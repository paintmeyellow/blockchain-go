package usecase

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"blockchain-go/blockchain"
)

type CreateBlockchainUcase struct {
	bc *blockchain.Blockchain
	tr trace.Tracer
}

func NewCreateBlockchainUcase(bc *blockchain.Blockchain) *CreateBlockchainUcase {
	return &CreateBlockchainUcase{
		bc: bc,
		tr: otel.Tracer("usecase"),
	}
}

func (ucase *CreateBlockchainUcase) Handle(ctx context.Context, addr string) error {
	ctx, span := ucase.tr.Start(ctx, "CreateBlockchainUcase.Handle")
	defer span.End()

	// func() {
	// 	_, span := ucase.tr.Start(ctx, "test")
	// 	time.Sleep(2 * time.Second)
	// 	span.End()
	// }()

	if err := ucase.bc.Create(ctx, addr); err != nil {
		err = fmt.Errorf("blockchain.create: %w", err)
		span.SetStatus(codes.Error, "operation failed")
		span.RecordError(err)
		return err
	}
	return nil
}
