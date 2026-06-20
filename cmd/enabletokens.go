/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"solana-cli/internal/util"

	"github.com/spf13/cobra"
)

// enabletokensCmd represents the enabletokens command
var enabletokensCmd = &cobra.Command{
	Use:   "enabletokens",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		data := util.GetTokensInfo()
		util.SaveTokensInfo(&data)
	},
}

func init() {
	rootCmd.AddCommand(enabletokensCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enabletokensCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enabletokensCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
