package util

import (
	"crypto/ed25519"
	"encoding/base64"
)

func b64decode(b64 string) []byte {
	r, _ := base64.StdEncoding.DecodeString(b64)
	return r
}

func GetMessage(tx Tx) []byte {
	raw := []byte{}
	raw = append(raw, tx.SignaturesLength...)
	for _, b := range tx.Signatures {
		raw = append(raw, b...)
	}
	raw = append(raw, tx.Version...)
	raw = append(raw, tx.SignerAccounts...)
	raw = append(raw, tx.ReadonlySignerAccounts...)
	raw = append(raw, tx.ReadonlyNonsignerAccounts...)
	raw = append(raw, tx.AllAccounts...)
	for _, b := range tx.WritableSigner {
		raw = append(raw, b...)
	}
	for _, b := range tx.ReadonlySigner {
		raw = append(raw, b...)
	}
	for _, b := range tx.WritableNonsigner {
		raw = append(raw, b...)
	}
	for _, b := range tx.ReadonlyNonsigner {
		raw = append(raw, b...)
	}
	raw = append(raw, tx.RecentBlockhash...)
	raw = append(raw, tx.InstructionsCount...)
	for _, ix := range tx.Instructions {
		raw = append(raw, ix.ProgramIdIndex...)
		raw = append(raw, ix.AccountsCount...)
		raw = append(raw, ix.Accounts...)
		raw = append(raw, ix.DataLength...)
		raw = append(raw, ix.Data...)
	}
	raw = append(raw, tx.AdtableLookup...)

	return raw
}

func GetSignature(msg []byte, priv []byte) []byte {
	privKey := ed25519.PrivateKey(priv)
	sign := ed25519.Sign(privKey, msg)
	return sign
}

func SignTx(tx Tx, priv []byte) Tx {
	msg := GetMessage(tx)
	signature := GetSignature(msg, priv)
	tx.SignaturesLength = tx.SignerAccounts
	tx.Signatures = append(tx.Signatures, signature)
	return tx
}

func GetCols(txtype string) []string {
	if txtype == "Payment" {
		return []string{
			"TransactionType",
			"Sequence",
			"DestinationTag",
			"Amount",
			"Fee",
			"SigningPubKey",
			"TxnSignature",
			"Account",
			"Destination",
		}
	} else if txtype == "AccountDelete" {
		return []string{
			"TransactionType",
			"Sequence",
			"DestinationTag",
			"Fee",
			"SigningPubKey",
			"TxnSignature",
			"Account",
			"Destination",
		}
	}
	return nil
}
