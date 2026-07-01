package util

import (
	"encoding/csv"
	"os"
)

func MkdirOrNothing(dir string) {
	os.MkdirAll(dir, 0755)
}

func LoadAccounts() [][]string {
	f, _ := os.Open(RelativeToAbsolute("ref", "SOL_keypair.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func LoadDestinations() [][]string {
	f, _ := os.Open(RelativeToAbsolute("ref", "SOL_destinations.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func SaveKeypair(name string, address string, priv string) {
	MkdirOrNothing("ref")
	f, _ := os.OpenFile("ref/SOL_keypair.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	f, _ := os.OpenFile("ref/SOL_destinations.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	row := []byte{}
	row = append(row, []byte(name)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, []byte("\n")...)
	f.Write(row)
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
