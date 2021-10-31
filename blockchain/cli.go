package blockchain

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

type cli struct {
	BC *Blockchain
}

func NewCLI(bc *Blockchain) *cli {
	return &cli{BC: bc}
}

func (cli *cli) Run(ctx context.Context) error {
	rootCmd := &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
		},
	}
	rootCmd.AddCommand(cli.addBlock())
	rootCmd.AddCommand(cli.printChain())
	return rootCmd.ExecuteContext(ctx)
}

func (cli *cli) addBlock() *cobra.Command {
	var data string
	cmd := cobra.Command{
		Use:   "addblock",
		Short: "Add block to blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			cli.BC.AddBlock(data)
			fmt.Println("Success!")
		},
	}
	cmd.PersistentFlags().StringVarP(&data, "data", "d", "", "Block data")
	return &cmd
}

func (cli *cli) printChain() *cobra.Command {
	return &cobra.Command{
		Use:   "printchain",
		Short: "Print blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			bci := cli.BC.Iterator()
			for {
				block := bci.Next()
				fmt.Printf("Prev: %x\n", block.PrevBlockHash)
				fmt.Printf("Data: %s\n", block.Data)
				fmt.Printf("Hash: %x\n", block.Hash)
				pow := NewProofOfWork(block)
				fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
				fmt.Println()
				if len(block.PrevBlockHash) == 0 {
					break
				}
			}
		},
	}
}
