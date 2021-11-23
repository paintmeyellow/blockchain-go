package blockchain

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

type cli struct {
	dbFile string
}

func NewCLI() *cli {
	return &cli{
		dbFile: "blockchain.db",
	}
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
	payto, err := cli.payto()
	if err != nil {
		return err
	}
	rootCmd.AddCommand(createBlockchainCmd)
	rootCmd.AddCommand(balanceCmd)
	rootCmd.AddCommand(payto)
	return rootCmd.ExecuteContext(ctx)
}

func (cli *cli) balance() (*cobra.Command, error) {
	var addr string
	cmd := cobra.Command{
		Use:   "balance",
		Short: "Get address balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			bc, err := NewBlockchain(cli.dbFile)
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

func (cli *cli) payto() (*cobra.Command, error) {
	var (
		from   string
		to     string
		amount int
	)
	cmd := cobra.Command{
		Use:   "payto",
		Short: "Pay to address",
		RunE: func(cmd *cobra.Command, args []string) error {
			bc, err := NewBlockchain(cli.dbFile)
			if err != nil {
				return err
			}
			defer bc.DB.Close()
			tx, err := NewTx(from, to, amount, bc)
			if err != nil {
				return err
			}
			if err = bc.MineBlock([]*Tx{tx}); err != nil {
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

func (cli *cli) createBlockchain() (*cobra.Command, error) {
	var addr string
	cmd := cobra.Command{
		Use:   "createblockchain",
		Short: "Create new blockchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			bc, err := CreateBlockchain(addr, cli.dbFile)
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
