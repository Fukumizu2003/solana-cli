/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"solana-cli/internal/util"

	"github.com/spf13/cobra"
)

var set_password string
var set_name string
var set_address string
var adonly bool

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if !adonly {
			privkey, pubkey := util.NewKeypair()
			address := util.PubkeyToAddress([]byte(pubkey))
			privkey_cr := util.AesEncrypt([]byte(privkey), []byte(set_password))
			util.SaveKeypair(set_name, address, base64.StdEncoding.EncodeToString(privkey_cr))
		} else {
			if set_address == "" {
				fmt.Println("-aフラグでアドレスを指定してください。")
				return
			}
			util.SaveAddress(set_name, set_address)
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&set_password, "password", "p", "", "")
	registerCmd.Flags().StringVarP(&set_name, "name", "n", "Default", "")
	registerCmd.Flags().StringVarP(&set_address, "address", "a", "", "")
	registerCmd.Flags().BoolVarP(&adonly, "addressonly", "o", false, "")
}
