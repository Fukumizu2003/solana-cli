/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var addrDebug string
var tokenDebug string

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)

	debugCmd.Flags().StringVarP(&addrDebug, "address", "a", "", "")
	debugCmd.Flags().StringVarP(&tokenDebug, "token", "t", "", "")
}
