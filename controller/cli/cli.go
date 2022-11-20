package cli

import (
	"context"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"blockchain-go/usecase"
)

type cli struct {
	getBalanceUcase       *usecase.GetBalanceUcase
	payToUcase            *usecase.PayToUcase
	createBlockchainUcase *usecase.CreateBlockchainUcase
	tr                    trace.Tracer
}

func New(
	getBalanceUcase *usecase.GetBalanceUcase,
	payToUcase *usecase.PayToUcase,
	createBlockchainUcase *usecase.CreateBlockchainUcase,
) *cli {
	return &cli{
		getBalanceUcase:       getBalanceUcase,
		payToUcase:            payToUcase,
		createBlockchainUcase: createBlockchainUcase,
		tr:                    otel.Tracer("cli"),
	}
}

func (cli *cli) Run(ctx context.Context) error {
	rootCmd := &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	createBlockchainCmd, err := cli.createChain(ctx)
	if err != nil {
		return err
	}
	balanceCmd, err := cli.balance(ctx)
	if err != nil {
		return err
	}
	payto, err := cli.payto(ctx)
	if err != nil {
		return err
	}
	rootCmd.AddCommand(createBlockchainCmd)
	rootCmd.AddCommand(balanceCmd)
	rootCmd.AddCommand(payto)
	return rootCmd.ExecuteContext(ctx)
}

func (cli *cli) balance(ctx context.Context) (*cobra.Command, error) {
	var addr string
	cmd := cobra.Command{
		Use:   "balance",
		Short: "Get address balance",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx, span := cli.tr.Start(ctx, "cli.balance")
			defer span.End()
			balance := cli.getBalanceUcase.Handle(ctx, addr)
			cmd.Printf("Balance of '%s': %d\n", addr, balance.Value)
		},
	}
	cmd.Flags().StringVarP(&addr, "addr", "", "", "Balance address")
	if err := cmd.MarkFlagRequired("addr"); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (cli *cli) payto(ctx context.Context) (*cobra.Command, error) {
	var (
		from   string
		to     string
		amount int
	)
	cmd := cobra.Command{
		Use:   "payto",
		Short: "Pay to address",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx, span := cli.tr.Start(ctx, "cli.payto")
			defer span.End()
			if err := cli.payToUcase.Handle(ctx, from, to, amount); err != nil {
				span.RecordError(err)
				cmd.PrintErrf("operation failed: %s", err)
				return
			}
			cmd.Println("Success!")
		},
	}
	cmd.Flags().StringVarP(&from, "from", "", "", "From address")
	cmd.Flags().StringVarP(&to, "to", "", "", "To address")
	cmd.Flags().IntVarP(&amount, "amount", "", 0, "Amount")
	if err := cmd.MarkFlagRequired("from"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("to"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("amount"); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (cli *cli) createChain(ctx context.Context) (*cobra.Command, error) {
	var addr string
	cmd := cobra.Command{
		Use:   "create-chain",
		Short: "Create new blockchain",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx, span := cli.tr.Start(ctx, "cli.createChain")
			defer span.End()
			if err := cli.createBlockchainUcase.Handle(ctx, addr); err != nil {
				span.SetStatus(codes.Error, "operation failed")
				span.RecordError(err)
				cmd.PrintErrf("operation failed: %s", err)
				return
			}
			cmd.Println("Success!")
		},
	}
	cmd.Flags().StringVarP(&addr, "addr", "", "", "Rewards address")
	if err := cmd.MarkFlagRequired("addr"); err != nil {
		return nil, err
	}
	return &cmd, nil
}
