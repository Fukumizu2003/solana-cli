/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"solana-cli/internal/config"

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
		st, er := config.SetAccount(mainAcName)
		if er != nil {
			fmt.Println(er)
			return
		}
		config.SaveConfig(*st)
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
