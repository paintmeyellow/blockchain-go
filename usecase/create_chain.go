package usecase

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := ucase.tr.Start(ctx, "usecase.create_chain")
	defer span.End()
	span.SetAttributes(attribute.String("addr", addr))

	err := ucase.bc.Create(ctx, addr)
	if err != nil {
		span.SetStatus(codes.Error, "operation failed")
		span.RecordError(fmt.Errorf("bc.Create: %w", err))
	}
	return err
}
