package config

import (
	"errors"
	"os"
	"solana-cli/internal/util"

	"github.com/joho/godotenv"
)

type State struct {
	Name    string
	Address string
	Key     string
}

func SetAccount(name string) (*State, error) {
	var state State

	accounts := util.LoadAccounts()
	flag := false
	for _, ac := range accounts {
		if name == ac[0] {
			state.Name = name
			state.Address = ac[1]
			state.Key = ac[2]
			flag = true
			break
		}
	}
	if !flag {
		return nil, errors.New("このアカウント名は存在しません。")
	}
	return &state, nil
}

func GetAccount() *State {
	godotenv.Load()
	var state State
	state.Name = os.Getenv("NAME_SOL")
	state.Address = os.Getenv("ADDRESS_SOL")
	state.Key = os.Getenv("PRIVKEY_ENCRYPTED_SOL")
	return &state
}

func SaveConfig(st State) {
	curr, err := godotenv.Read(".env")
	if err != nil {
		curr = make(map[string]string)
	}
	curr["NAME_SOL"] = st.Name
	curr["ADDRESS_SOL"] = st.Address
	curr["PRIVKEY_ENCRYPTED_SOL"] = st.Key
	godotenv.Write(curr, ".env")
}
