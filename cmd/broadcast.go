/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"solana-cli/internal/util"

	"fmt"

	"github.com/spf13/cobra"
)

var broadcastCmd = &cobra.Command{
	Use:   "broadcast",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tx := util.LoadTx()
		res := util.Broadcast(tx)
		if res["error"] != nil {
			errfield := res["error"].(map[string]interface{})
			fmt.Println("ERROR: " + errfield["message"].(string))
		} else if res["result"] != nil {
			result := res["result"]
			fmt.Println("SUCCEED: " + result.(string))
		}
	},
}

func init() {
	rootCmd.AddCommand(broadcastCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// broadcastCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// broadcastCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
