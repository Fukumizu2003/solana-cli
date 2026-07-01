package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"filippo.io/edwards25519"
)

func PubkeyToAddress(pub []byte) string {
	address := B58Encode(pub)
	return address
}

func AddressToPubkey(addr string) []byte {
	pubk := B58Decode(addr)
	return pubk
}

func isED25519(pubKeyBytes []byte) bool {
	if len(pubKeyBytes) != 32 {
		return false
	}
	var p edwards25519.Point
	_, err := p.SetBytes(pubKeyBytes)
	return err == nil
}

func GetTokenInfoByMint(mint string) map[string]interface{} {
	var tokendata map[string]interface{}
	url := "https://api.mainnet-beta.solana.com"
	payload := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getAccountInfo",
		"params": [
			"` + mint + `",
			{
				"encoding": "jsonParsed"
			}
		]
	}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&tokendata); err != nil {
		panic(err)
	}
	return tokendata
}

func ReadOwner(info map[string]interface{}) (string, error) {
	if info["result"] == nil {
		return "", errors.New("\"result\"カラムなし")
	}
	result := info["result"].(map[string]interface{})
	if result["value"] == nil {
		return "", errors.New("\"value\"カラムなし")
	}
	value := result["value"].(map[string]interface{})
	if value["owner"] == nil {
		return "", errors.New("\"owner\"カラムなし")
	}
	return value["owner"].(string), nil
}

func GetTokenProgramATA(wallet string, token string) (string, error) {
	tokenInfo := GetTokenInfoByMint(token)
	program, e := ReadOwner(tokenInfo)
	if e != nil {
		return "", e
	}
	ata := CalcATA(wallet, token, program)
	return ata, nil
}

func CalcATA(wallet string, token string, tokenprogram string) string {
	walletBytes := B58Decode(wallet)
	tokenProgramBytes := B58Decode(tokenprogram)
	mintBytes := B58Decode(token)
	ataProgramBytes := B58Decode("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")

	pdaSalt := []byte("ProgramDerivedAddress")

	bumpB := byte(255)

	for {
		hasher := sha256.New()
		hasher.Write(walletBytes)
		hasher.Write(tokenProgramBytes)
		hasher.Write(mintBytes)
		hasher.Write([]byte{bumpB})
		hasher.Write(ataProgramBytes)
		hasher.Write(pdaSalt)

		ataB := hasher.Sum(nil)

		if !isED25519(ataB) {
			return B58Encode(ataB)
		}

		bumpB--
	}
}

func GetAddressTokens(addr string) ([]map[string]interface{}, error) {
	var tokendatas []map[string]interface{}
	url := "https://api.mainnet-beta.solana.com"
	payload := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getTokenAccountsByOwner",
		"params": [
			"` + addr + `",
			{ "programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA" },
			{ "encoding": "jsonParsed" }
		]
	}`)
	payload2 := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getTokenAccountsByOwner",
		"params": [
			"` + addr + `",
			{ "programId": "TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb" },
			{ "encoding": "jsonParsed" }
		]
	}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	resp2, err := http.Post(url, "application/json", bytes.NewBuffer(payload2))
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	defer resp.Body.Close()
	var result map[string]interface{}
	var result2 map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	if err := json.NewDecoder(resp2.Body).Decode(&result2); err != nil {
		panic(err)
	}
	rescol := result["result"].(map[string]interface{})
	rescol2 := result2["result"].(map[string]interface{})
	valcol := rescol["value"].([]interface{})
	valcol2 := rescol2["value"].([]interface{})
	for _, tkdata := range append(valcol, valcol2...) {
		tk, ok := tkdata.(map[string]interface{})
		if !ok {
			continue
		}
		tokendatas = append(tokendatas, tk["account"].(map[string]interface{}))
	}
	return tokendatas, nil
}

func ReadTokensInfo(datas []map[string]interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	for _, data := range datas {
		mp := make(map[string]interface{})
		datacl := data["data"].(map[string]interface{})
		owner := data["owner"].(string)
		parsedcl := datacl["parsed"].(map[string]interface{})
		infocl := parsedcl["info"].(map[string]interface{})
		mintad := infocl["mint"].(string)
		blcl := infocl["tokenAmount"].(map[string]interface{})
		bl, _ := strconv.Atoi(blcl["amount"].(string))
		dc := int(blcl["decimals"].(float64))
		blStr := blcl["uiAmountString"].(string)
		ti := readTokenInfo(mintad)
		name := ti.name
		sym := ti.symbol
		mp["name"] = name
		mp["symbol"] = sym
		mp["decimal"] = dc
		mp["mint"] = mintad
		mp["owner"] = owner
		mp["balance"] = bl
		mp["balancestr"] = blStr
		res = append(res, mp)
	}
	return res
}

type tokenInfo struct {
	name   string
	symbol string
	id     string
	owner  string
}

func readTokenInfo(token string) tokenInfo {
	var ti tokenInfo
	ti.id = token
	tkInfo := LoadTokensInfo()
	for _, tkdata := range tkInfo {
		if tkdata["id"] == nil {
			continue
		}
		if token == tkdata["id"].(string) {
			ti.symbol = tkdata["symbol"].(string)
			ti.name = tkdata["name"].(string)
			ti.owner = tkdata["tokenProgram"].(string)
			break
		}
	}
	return ti
}

func GetTokensInfo() []map[string]interface{} {
	url := "https://api.jup.ag/tokens/v2/tag?query=verified"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result []map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	return result
}

func SaveTokensInfo(data *[]map[string]interface{}) {
	save, _ := json.MarshalIndent(*data, "", "    ")
	os.WriteFile(RelativeToAbsolute("ref", "SOL_TokenReference.json"), save, 0644)
}

func LoadTokensInfo() []map[string]interface{} {
	data, _ := os.ReadFile(RelativeToAbsolute("ref", "SOL_TokenReference.json"))
	var res []map[string]interface{}
	json.Unmarshal(data, &res)
	return res
}

func GetATAInfo(addr string, mint string) (map[string]interface{}, error) {
	var res map[string]interface{}
	ata, e := GetTokenProgramATA(addr, mint)
	if e != nil {
		return nil, e
	}
	url := "https://api.mainnet-beta.solana.com"
	payload := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getAccountInfo",
		"params": [
			"` + ata + `",
			{
				"encoding": "jsonParsed"
			}
		]
	}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		panic(err)
	}
	return res, nil
}

func ExistATA(info map[string]interface{}) (bool, error) {
	if info["result"] == nil {
		return false, errors.New("情報取得失敗")
	}
	result := info["result"].(map[string]interface{})
	if result["value"] == nil {
		return false, nil
	}
	return true, nil
}
