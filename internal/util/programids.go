package util

import (
	"encoding/hex"
)

func GetProgramId(kind string) []byte {
	r := []byte{}
	if kind == "System" {
		r, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	}
	return r
}
