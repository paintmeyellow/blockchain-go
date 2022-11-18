package main

import (
	"context"
	"log"

	"github.com/boltdb/bolt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"blockchain-go/blockchain"
	"blockchain-go/controller/cli"
	"blockchain-go/usecase"
)

const (
	jaegerCollectorURL = "http://localhost:14268/api/traces"
	serviceName        = "blockchain-demo"
	serviceVersion     = "v0.0.1"
	environment        = "development"
)

func main() {
	// create the jaeger exporter
	jaegerExp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(jaegerCollectorURL),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(jaegerExp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			attribute.String("environment", environment),
		)),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalln(err)
		}
	}()

	// set tracer provider as singletone
	otel.SetTracerProvider(tp)

	// setup boltdb
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// setup dependencies
	bc, err := blockchain.Open(db)
	if err != nil {
		log.Fatalln(err)
	}
	// create and start CLI
	client := cli.New(
		usecase.NewGetBalanceUcase(bc),
		usecase.NewPayToUcase(bc),
		usecase.NewCreateBlockchainUcase(db),
	)
	if err := client.Run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
