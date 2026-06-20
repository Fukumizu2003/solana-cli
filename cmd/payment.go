/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"solana-cli/internal/config"
	"solana-cli/internal/util"

	"github.com/spf13/cobra"
)

var sendAmountStr string
var sendDestination string
var sendAll bool
var sendToken string

// paymentCmd represents the payment command
var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if sendAmountStr == "" && !sendAll {
			fmt.Println("--all（全額送金）もしくは-aフラグで送金額を指定してください。")
			return
		}
		if sendAmountStr != "" && sendAll {
			fmt.Println("--allと送金額を同時に指定しないでください。")
			return
		}
		if sendDestination == "" {
			fmt.Println("-dで送金先を指定してください。")
			return
		}
		tx := util.Tx{}
		blockHash := util.GetBlockHash()
		address := config.GetAccount().Address
		senderPubkey := util.AddressToPubkey(address)
		destAddress := util.NameToAddress(sendDestination)
		destPubkey := util.AddressToPubkey(destAddress)
		if sendToken == "" {
			var lamps int
			if !sendAll {
				lamps = util.SolToLamps(sendAmountStr)
			} else {
				lamps = util.GetBalance(address) - 5000
			}
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
			ix := util.Ix{}
			ix.ProgramIdIndex = []byte{0x02}
			ix.AccountsCount = []byte{0x02}
			ix.Accounts = []byte{0x00, 0x01}
			ix.DataLength = []byte{0x0c}
			lampBytes := util.IntToBytesFixed(lamps, 8)
			ix.Data = append([]byte{0x02, 0x00, 0x00, 0x00}, lampBytes...)
			tx.Instructions = []util.Ix{ix}
		} else {
			tokens, e := util.GetAddressTokens(address)
			if e != nil {
				fmt.Println(e)
				return
			}
			infos := util.ReadTokensInfo(tokens)
			var targetToken map[string]interface{}
			exist := false
			for _, info := range infos {
				if strings.EqualFold(info["name"].(string), sendToken) || strings.EqualFold(info["symbol"].(string), sendToken) || strings.EqualFold(info["mint"].(string), sendToken) {
					targetToken = info
					exist = true
					break
				}
			}
			if !exist {
				fmt.Println("指定のトークンを保有していません。")
				return
			}
			mint := targetToken["mint"].(string)
			dec := targetToken["decimal"].(int)
			tokenProgramId := targetToken["owner"].(string)
			am := util.FloatStrToInt(sendAmountStr, dec)
			ataInfo, _ := util.GetATAInfo(destAddress, mint)
			fromAta, _ := util.GetTokenProgramATA(address, mint)
			ataExist, er := util.ExistATA(ataInfo)
			if er != nil {
				fmt.Println(er)
				return
			}
			toAta := util.CalcATA(destAddress, mint, tokenProgramId)
			if er != nil {
				fmt.Println(er)
				return
			}
			if !ataExist {
				tx.Version = []byte{0x80}
				tx.SignerAccounts = []byte{0x01}
				tx.ReadonlySignerAccounts = []byte{0x00}
				tx.ReadonlyNonsignerAccounts = []byte{0x06}
				tx.AllAccounts = []byte{0x09}
				tx.WritableSigner = [][]byte{util.B58Decode(address)}
				tx.WritableNonsigner = [][]byte{util.B58Decode(fromAta), util.B58Decode(toAta)}
				tx.ReadonlyNonsigner = [][]byte{
					util.B58Decode(destAddress),
					util.B58Decode(mint),
					util.GetProgramId("System"),
					util.B58Decode(tokenProgramId),
					util.GetProgramId("Rent Program"),
					util.GetProgramId("Associated Token Program"),
				}
				tx.InstructionsCount = []byte{0x02}
				ix1 := util.Ix{}
				ix2 := util.Ix{}
				ix1.ProgramIdIndex = []byte{0x08}
				ix1.AccountsCount = []byte{0x07}
				ix1.Accounts = []byte{0x00, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
				ix1.DataLength = []byte{0x00}
				ix2.ProgramIdIndex = []byte{0x06}
				ix2.AccountsCount = []byte{0x03}
				ix2.Accounts = []byte{0x01, 0x02, 0x00}
				ix2.DataLength = []byte{0x09}
				amb := util.IntToBytesFixed(am, 8)
				ix2.Data = append([]byte{0x03}, amb...)
				tx.Instructions = []util.Ix{ix1, ix2}
				tx.AdtableLookup = []byte{0x00}
				tx.RecentBlockhash = util.AddressToPubkey(blockHash)
			} else {
				tx.Version = []byte{0x80}
				tx.SignerAccounts = []byte{0x01}
				tx.ReadonlySignerAccounts = []byte{0x00}
				tx.ReadonlyNonsignerAccounts = []byte{0x01}
				tx.AllAccounts = []byte{0x04}
				tx.WritableSigner = [][]byte{util.B58Decode(address)}
				tx.WritableNonsigner = [][]byte{util.B58Decode(fromAta), util.B58Decode(toAta)}
				tx.ReadonlyNonsigner = [][]byte{util.B58Decode(tokenProgramId)}
				tx.InstructionsCount = []byte{0x01}
				ix := util.Ix{}
				ix.ProgramIdIndex = []byte{0x03}
				ix.AccountsCount = []byte{0x03}
				ix.Accounts = []byte{0x01, 0x02, 0x00}
				ix.DataLength = []byte{0x09}
				amb := util.IntToBytesFixed(am, 8)
				ix.Data = append([]byte{0x03}, amb...)
				tx.Instructions = []util.Ix{ix}
				tx.AdtableLookup = []byte{0x00}
				tx.RecentBlockhash = util.AddressToPubkey(blockHash)
			}
		}
		util.SaveTx(tx)
	},
}

func init() {
	rootCmd.AddCommand(paymentCmd)

	paymentCmd.Flags().StringVarP(&sendAmountStr, "amount", "a", "", "")
	paymentCmd.Flags().StringVarP(&sendDestination, "destination", "d", "", "")
	paymentCmd.Flags().StringVar(&sendToken, "token", "", "")
	paymentCmd.Flags().BoolVar(&sendAll, "all", false, "To send maximum.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// paymentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// paymentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
