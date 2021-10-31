package blockchain

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

type cli struct {
}

func NewCLI() *cli {
	return &cli{}
}

func (cli *cli) Run(ctx context.Context) error {
	rootCmd := &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	createBlockchainCmd, err := cli.createBlockchain()
	if err != nil {
		return err
	}
	balanceCmd, err := cli.balance()
	if err != nil {
		return err
	}
	rootCmd.AddCommand(createBlockchainCmd)
	rootCmd.AddCommand(balanceCmd)
	return rootCmd.ExecuteContext(ctx)
}

func (cli *cli) balance() (*cobra.Command, error) {
	var addr string
	cmd := cobra.Command{
		Use:   "balance",
		Short: "Get address balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			bc, err := NewBlockchain()
			if err != nil {
				return err
			}
			defer bc.DB.Close()
			var balance int
			utxo := bc.UTXO(addr)
			for _, out := range utxo {
				balance += out.Value
			}
			fmt.Printf("Balance of '%s': %d\n", addr, balance)
			return nil
		},
	}
	cmd.Flags().StringVarP(&addr, "addr", "", "", "Balance address")
	if err := cmd.MarkFlagRequired("addr"); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (cli *cli) createBlockchain() (*cobra.Command, error) {
	var addr string
	cmd := cobra.Command{
		Use:   "createblockchain",
		Short: "Create new blockchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			bc, err := CreateBlockchain(addr)
			if err != nil {
				return err
			}
			defer bc.DB.Close()
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
