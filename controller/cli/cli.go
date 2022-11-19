package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"

	"blockchain-go/usecase"
)

const tracerName = "cli"

type cli struct {
	getBalanceUcase       *usecase.GetBalanceUcase
	payToUcase            *usecase.PayToUcase
	createBlockchainUcase *usecase.CreateBlockchainUcase
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
	}
}

func (cli *cli) Run(ctx context.Context) error {
	rootCmd := &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	createBlockchainCmd, err := cli.createBlockchain(ctx)
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
	ctx, span := otel.Tracer(tracerName).Start(ctx, "balance")
	defer span.End()

	var addr string
	cmd := cobra.Command{
		Use:   "balance",
		Short: "Get address balance",
		RunE: func(_ *cobra.Command, _ []string) error {
			balance := cli.getBalanceUcase.Handle(ctx, addr)
			fmt.Printf("Balance of '%s': %d\n", addr, balance.Value)
			return nil
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
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			err := cli.payToUcase.Handle(ctx, from, to, amount)
			if err != nil {
				return err
			}
			fmt.Println("Success!")
			return nil
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

func (cli *cli) createBlockchain(ctx context.Context) (*cobra.Command, error) {
	var addr string
	cmd := cobra.Command{
		Use:   "createblockchain",
		Short: "Create new blockchain",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := cli.createBlockchainUcase.Handle(addr)
			if err != nil {
				return err
			}
			fmt.Println("Done!")
			return nil
		},
	}
	cmd.Flags().StringVarP(&addr, "addr", "", "", "Rewards address")
	if err := cmd.MarkFlagRequired("addr"); err != nil {
		return nil, err
	}
	return &cmd, nil
}
