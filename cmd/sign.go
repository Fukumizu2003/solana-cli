/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"solana-cli/internal/util"

	"github.com/spf13/cobra"
)

var pw string

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, _, privcr := util.GetAccount()
		privkey, _ := util.AesDecrypt(util.B64Decode(privcr), []byte(pw))
		tx := util.LoadTx()
		signedtx := util.SignTx(tx, privkey)
		util.SaveTx(signedtx)
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&pw, "password", "p", "", "")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
