package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Tx struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

const subsidy = 50

func NewCoinbaseTx(to, data string) (*Tx, error) {
	if data == "" {
		data = fmt.Sprintf("reward to '%s'", to)
	}
	txin := TxInput{
		TxID:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := TxOutput{Value: subsidy, ScriptPubKey: to}
	tx := Tx{
		ID:   nil,
		Vin:  []TxInput{txin},
		Vout: []TxOutput{txout},
	}
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

func (tx Tx) Serialize() ([]byte, error) {
	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(tx); err != nil {
		return nil, err
	}
	return encoded.Bytes(), nil
}

func (tx *Tx) Hash() ([]byte, error) {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	data, err := txCopy.Serialize()
	if err != nil {
		return nil, err
	}
	hash = sha256.Sum256(data)
	return hash[:], nil
}
