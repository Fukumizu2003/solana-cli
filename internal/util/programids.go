package util

func GetProgramId(kind string) []byte {
	r := []byte{}
	switch kind {
	case "System":
		r = B58Decode("11111111111111111111111111111111")
	case "Token Program":
		r = B58Decode("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	case "Token 2022 Program":
		r = B58Decode("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")
	case "Associated Token Program":
		r = B58Decode("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
	case "Rent Program":
		r = B58Decode("SysvarRent111111111111111111111111111111111")
	}
	return r
}
