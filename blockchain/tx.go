package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

type Tx struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

const reward = 50

func NewTx(from, to string, amount int, acc int, utxo map[string][]int) (*Tx, error) {
	var inputs []TxInput
	var outputs []TxOutput
	for txid, outs := range utxo {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			return nil, err
		}
		for _, outIdx := range outs {
			input := TxInput{TxID: txID, Vout: outIdx, ScriptSig: from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{Value: amount, ScriptPubKey: to})
	if acc > amount {
		//send change back
		outputs = append(outputs, TxOutput{Value: acc - amount, ScriptPubKey: from})
	}
	tx := Tx{Vin: inputs, Vout: outputs}
	hash, err := tx.Hash()
	if err != nil {
		return nil, err
	}
	tx.ID = hash
	return &tx, nil
}

func NewCoinbaseTx(to, data string) (*Tx, error) {
	if data == "" {
		data = fmt.Sprintf("reward to '%s'", to)
	}
	txin := TxInput{
		TxID:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := TxOutput{Value: reward, ScriptPubKey: to}
	tx := Tx{Vin: []TxInput{txin}, Vout: []TxOutput{txout}}
	hash, err := tx.Hash()
	if err != nil {
		return nil, err
	}
	tx.ID = hash
	return &tx, nil
}

func (tx *Tx) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Tx) Serialize() ([]byte, error) {
	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(tx); err != nil {
		return nil, err
	}
	return encoded.Bytes(), nil
}

func (tx Tx) Hash() ([]byte, error) {
	var hash [32]byte
	tx.ID = []byte{}
	data, err := tx.Serialize()
	if err != nil {
		return nil, err
	}
	hash = sha256.Sum256(data)
	return hash[:], nil
}
