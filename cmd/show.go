/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"solana-cli/internal/config"
	"solana-cli/internal/util"
	"strings"

	"github.com/spf13/cobra"
)

var show_address bool
var show_balance bool
var show_all bool
var show_token string

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		st := config.GetAccount()
		accname := st.Name
		addr := st.Address
		if !show_address && !show_balance {
			fmt.Println("表示する内容を指定してください：")
			fmt.Println("-a：アドレス")
			fmt.Println("-b：残高")
		}
		if show_address {
			if !show_all {
				fmt.Println("Name:    " + accname)
				fmt.Println("Address: " + addr)
				util.ShowQRCode(addr)
			} else {
				accts := append(util.LoadAccounts(), util.LoadDestinations()...)
				mxlen := 0
				for _, acct := range accts {
					if len(acct[0]) > mxlen {
						mxlen = len(acct[0])
					}
				}
				fmt.Printf("%-*s |Address\n", mxlen, "Name")
				fmt.Println(strings.Repeat("-", mxlen+48))
				for _, acct := range accts {
					fmt.Printf("%-*s |%s\n", mxlen, acct[0], acct[1])
				}
			}
		}
		if show_balance {
			if show_token == "" {
				bal := util.GetBalance(addr)
				balsol := util.LampsToSol(util.IntToStr(bal))
				fmt.Println("Balance: " + balsol + " SOL")
				return
			}
			tokens, e := util.GetAddressTokens(addr)
			if e != nil {
				fmt.Println(e)
				return
			}
			infos := util.ReadTokensInfo(tokens)
			for _, info := range infos {
				if strings.EqualFold(info["name"].(string), show_token) || strings.EqualFold(info["symbol"].(string), show_token) || strings.EqualFold(info["mint"].(string), show_token) {
					fmt.Println("Balance: " + info["balancestr"].(string) + " " + info["symbol"].(string))
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().BoolVarP(&show_address, "address", "a", false, "")
	showCmd.Flags().BoolVarP(&show_balance, "balance", "b", false, "")
	showCmd.Flags().BoolVar(&show_all, "all", false, "")
	showCmd.Flags().StringVar(&show_token, "token", "", "")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
