package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func GetBlockHash() string {
	url := "https://api.mainnet-beta.solana.com"

	payload := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getLatestBlockhash",
		"params": [{"commitment": "finalized"}]
	}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			Value struct {
				Blockhash string `json:"blockhash"`
			} `json:"value"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	return result.Result.Value.Blockhash
}

func Broadcast(tx Tx) map[string]interface{} {
	url := "https://api.mainnet-beta.solana.com"
	txBytes := GetMessage(tx)
	txB64 := B64Encode(txBytes)

	payload := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "sendTransaction",
		"params": [
			"` + txB64 + `",
			{ "encoding": "base64" }
		]
	}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	MkdirOrNothing("temp")
	f, _ := os.OpenFile("temp/SOL_resp.json", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	save, _ := json.MarshalIndent(result, "", "    ")
	f.Write(save)
	return result
}

func GetBalance(addr string) int {
	url := "https://api.mainnet-beta.solana.com"

	payload := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getBalance",
		"params": [
			"` + addr + `",
			{"commitment": "finalized"}
		]
	}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	rescol := result["result"].(map[string]interface{})
	balfl := rescol["value"].(float64)
	return int(balfl)
}
