/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"solana-cli/internal/util"

	"github.com/spf13/cobra"
)

var sendAmountStr string
var sendDestination string
var sendAll bool

// paymentCmd represents the payment command
var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if sendAmountStr == "" && !sendAll {
			fmt.Println("--allフラグを立てるもしくは-aフラグで送金額を指定してください。")
			return
		}
		if sendAmountStr != "" && sendAll {
			fmt.Println("--allフラグと送金額を同時に指定しないでください。")
			return
		}
		if sendDestination == "" {
			fmt.Println("-dフラグで送金先を指定してください。")
			return
		}
		_, address, _ := util.GetAccount()
		senderPubkey := util.AddressToPubkey(address)
		destPubkey := util.AddressToPubkey(util.NameToAddress(sendDestination))
		var lamps int
		if !sendAll {
			lamps = util.SolToLamps(sendAmountStr)
		} else {
			lamps = util.GetBalance(address) - 5000
		}

		blockHash := util.GetBlockHash()

		tx := util.Tx{}
		ix := util.Ix{}
		tx.Version = []byte{0x80}
		tx.SignerAccounts = []byte{0x01}
		tx.ReadonlySignerAccounts = []byte{0x00}
		tx.ReadonlyNonsignerAccounts = []byte{0x01}
		tx.AllAccounts = []byte{0x03}
		tx.WritableSigner = [][]byte{senderPubkey}
		tx.WritableNonsigner = [][]byte{destPubkey}
		tx.ReadonlyNonsigner = [][]byte{util.GetProgramId("System")}
		tx.InstructionsCount = []byte{0x01}
		tx.AdtableLookup = []byte{0x00}
		tx.RecentBlockhash = util.AddressToPubkey(blockHash)
		ix.ProgramIdIndex = []byte{0x02}
		ix.AccountsCount = []byte{0x02}
		ix.Accounts = []byte{0x00, 0x01}
		ix.DataLength = []byte{0x0c}
		lampBytes, _ := util.IntToBytes(lamps, 8)
		ix.Data = append([]byte{0x02, 0x00, 0x00, 0x00}, lampBytes...)
		tx.Instructions = []util.Ix{ix}
		util.SaveTx(tx)
	},
}

func init() {
	rootCmd.AddCommand(paymentCmd)

	paymentCmd.Flags().StringVarP(&sendAmountStr, "amount", "a", "", "")
	paymentCmd.Flags().StringVarP(&sendDestination, "destination", "d", "", "")
	paymentCmd.Flags().BoolVar(&sendAll, "all", false, "To send maximum.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// paymentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// paymentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
