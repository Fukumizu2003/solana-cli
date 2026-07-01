package util

import (
	"encoding/json"
	"os"
)

type Tx struct {
	SignaturesLength          []byte
	Signatures                [][]byte
	Version                   []byte
	SignerAccounts            []byte
	ReadonlySignerAccounts    []byte
	ReadonlyNonsignerAccounts []byte
	AllAccounts               []byte
	WritableSigner            [][]byte
	ReadonlySigner            [][]byte
	WritableNonsigner         [][]byte
	ReadonlyNonsigner         [][]byte
	RecentBlockhash           []byte
	InstructionsCount         []byte
	Instructions              []Ix
	AdtableLookup             []byte
}

type Ix struct {
	ProgramIdIndex []byte
	AccountsCount  []byte
	Accounts       []byte
	DataLength     []byte
	Data           []byte
}

func LoadTx() Tx {
	var tx Tx
	f, _ := os.ReadFile(RelativeToAbsolute("temp", "SOL_transaction.json"))
	json.Unmarshal(f, &tx)
	return tx
}
