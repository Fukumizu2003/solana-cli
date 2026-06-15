package util

import (
	"encoding/csv"
	"encoding/json"
	"os"
)

func MkdirOrNothing(dir string) {
	os.MkdirAll(dir, 0755)
}

func LoadAccounts() [][]string {
	f, _ := os.Open(RelativeToAbsolute("ref", "keypair.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func LoadDestinations() [][]string {
	f, _ := os.Open(RelativeToAbsolute("ref", "destinations.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func SaveKeypair(name string, address string, priv string) {
	MkdirOrNothing("ref")
	f, _ := os.OpenFile("ref/keypair.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	row := []byte{}
	row = append(row, []byte(name)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, byte(','))
	row = append(row, priv...)
	row = append(row, []byte("\n")...)
	f.Write(row)
}

func SaveAddress(name string, address string) {
	MkdirOrNothing("ref")
	f, _ := os.OpenFile("ref/destinations.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	row := []byte{}
	row = append(row, []byte(name)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, []byte("\n")...)
	f.Write(row)
}

func SetAccount(name string, addr string, key string) {
	MkdirOrNothing("ref")
	f, _ := os.OpenFile("ref/mainaccount.json", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	datamap := map[string]string{
		"Name":    name,
		"Address": addr,
		"Key":     key,
	}
	jsondata, _ := json.MarshalIndent(datamap, "", "    ")
	f.Write(jsondata)
}

func GetAccount() (string, string, string) {
	f, _ := os.ReadFile(RelativeToAbsolute("ref", "mainaccount.json"))
	var receipt map[string]string
	json.Unmarshal(f, &receipt)
	return receipt["Name"], receipt["Address"], receipt["Key"]
}

func NameToAddress(name string) string {
	accts := LoadAccounts()
	daccts := LoadDestinations()
	for _, ac := range accts {
		if ac[0] == name {
			return ac[1]
		}
	}
	for _, ac := range daccts {
		if ac[0] == name {
			return ac[1]
		}
	}
	return name
}
