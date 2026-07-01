/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"solana-cli/internal/util"

	"github.com/spf13/cobra"
)

var delete_name string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if delete_name == "" {
			return fmt.Errorf("削除するアカウント名を-nフラグで指定してください。")
		}
		accounts := util.LoadAccounts()
		destinations := util.LoadDestinations()
		new_accounts := [][]string{}
		new_dests := [][]string{}
		acflag := false
		deflag := false
		for _, ac := range accounts {
			if ac[0] != delete_name {
				new_accounts = append(new_accounts, ac)
				acflag = true
			}
		}
		for _, ac := range destinations {
			if ac[0] != delete_name {
				new_dests = append(new_dests, ac)
				deflag = true
			}
		}
		if !acflag && !deflag {
			return fmt.Errorf("アカウント名が存在しません。")
		}
		if acflag {
			var buf bytes.Buffer
			writer := csv.NewWriter(&buf)
			writer.WriteAll(new_accounts)
			os.WriteFile(util.RelativeToAbsolute("ref", "SOL_keypair.csv"), buf.Bytes(), 0644)
		}
		if deflag {
			var buf bytes.Buffer
			writer := csv.NewWriter(&buf)
			writer.WriteAll(new_dests)
			os.WriteFile(util.RelativeToAbsolute("ref", "SOL_destinations.csv"), buf.Bytes(), 0644)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&delete_name, "name", "n", "", "")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
