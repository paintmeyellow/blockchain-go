package blockchain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestBlockchain_CreateBlockchain(t *testing.T) {
	var (
		path          = "/tmp/test_blockchain.db"
		rewardAddress = "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	)

	t.Run("created", func(t *testing.T) {
		RequirePruneBlockchain(t, path)
		_, err := CreateBlockchain(rewardAddress, path)
		assert.NoError(t, err)
	})

	t.Run("already_exists", func(t *testing.T) {
		RequirePruneBlockchain(t, path)
		_, err := CreateBlockchain(rewardAddress, path)
		require.NoError(t, err)
		_, err = CreateBlockchain(rewardAddress, path)
		assert.ErrorIs(t, err, ErrBlockchainAlreadyExists)
	})

	t.Run("has_genesis_block", func(t *testing.T) {
		RequirePruneBlockchain(t, path)
		bc, err := CreateBlockchain(rewardAddress, path)
		require.NoError(t, err)
		assert.NotEmpty(t, bc.tip)
		bci := bc.Iterator()
		var b *Block
		var blocksCount int
		for {
			b = bci.Next()
			blocksCount++
			if len(b.PrevBlockHash) == 0 {
				break
			}
		}
		assert.Equal(t, 1, blocksCount)
	})

	t.Run("address_recieved_reward", func(t *testing.T) {
		RequirePruneBlockchain(t, path)
		bc, err := CreateBlockchain(rewardAddress, path)
		require.NoError(t, err)
		utxo := bc.UTXO(rewardAddress)
		assert.Equal(t, reward, UTXOBalance(utxo))
	})
}

func TestBlockchain_PayTo(t *testing.T) {
	var (
		rewardAddress = "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	)
	t.Run("create_transaction", func(t *testing.T) {
		var (
			from   = rewardAddress
			to     = "12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX"
			amount = 5
		)
		bc := RequireInitBlockchain(t, rewardAddress)
		_, err := NewTx(from, to, amount, bc)
		assert.NoError(t, err)
	})

	t.Run("not_enough_funds", func(t *testing.T) {
		var (
			from   = "bc1qex0aqq8mxqfh4cpl62eg755836djjx20yzuuu8"
			to     = ""
			amount = 1
		)
		bc := RequireInitBlockchain(t, rewardAddress)
		_, err := NewTx(from, to, amount, bc)
		assert.ErrorIs(t, err, ErrNotEnoughFunds)
	})

	t.Run("correct_sender_and_reciever_balances", func(t *testing.T) {
		var (
			from   = rewardAddress
			to     = "12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX"
			amount = 5
		)
		bc := RequireInitBlockchain(t, rewardAddress)
		fromBalance := UTXOBalance(bc.UTXO(from))
		tx := RequireNewTx(t, from, to, amount, bc)
		require.NoError(t, bc.MineBlock([]*Tx{tx}))
		fromBalance = fromBalance - amount
		toBalance := amount
		assert.Equal(t, fromBalance, UTXOBalance(bc.UTXO(from)))
		assert.Equal(t, toBalance, UTXOBalance(bc.UTXO(to)))
	})
}

func UTXOBalance(utxo []TxOutput) int {
	var balance int
	for _, out := range utxo {
		balance += out.Value
	}
	return balance
}

func RequirePruneBlockchain(t *testing.T, path string) {
	if _, err := os.Stat(path); err == nil {
		require.NoError(t, os.Remove(path))
	}
}

func RequireNewTx(t *testing.T, from, to string, amount int, bc *Blockchain) *Tx {
	tx, err := NewTx(from, to, amount, bc)
	require.NoError(t, err)
	return tx
}

func RequireInitBlockchain(t *testing.T, addr string) *Blockchain {
	path := "/tmp/test_blockchain.db"
	RequirePruneBlockchain(t, path)
	bc, err := CreateBlockchain(addr, path)
	require.NoError(t, err)
	return bc
}
