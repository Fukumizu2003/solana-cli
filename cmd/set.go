/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"solana-cli/internal/util"

	"github.com/spf13/cobra"
)

var mainAcName string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if mainAcName == "" {
			fmt.Println("-nフラグによりアカウント名を指定してください。")
			return
		}
		accts := util.LoadAccounts()
		addr := ""
		key := ""
		for _, ac := range accts {
			if ac[0] == mainAcName {
				addr = ac[1]
				key = ac[2]
				break
			}
		}
		if addr == "" {
			fmt.Println("指定のアカウント名は存在しません。")
			return
		}
		util.SetAccount(mainAcName, addr, key)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&mainAcName, "name", "n", "", "")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
