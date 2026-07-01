package util

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

func B64Encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	return encoded
}

func B64Decode(msg string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(msg)
	return decoded
}

func BytesToInt(orig []byte) int {
	copy := orig[:]
	sum := 0
	scale := 1
	for _, by := range copy {
		sum += int(by) * scale
		scale <<= 8
	}
	return sum
}

func IntToBytes(orig int) []byte {
	bytes := []byte{}
	for orig != 0 {
		bytes = append(bytes, byte(orig&0xff))
		orig = orig >> 8
	}
	return bytes
}

func IntToBytesFixed(orig int, leng int) []byte {
	bytes := []byte{}
	for orig != 0 {
		bytes = append(bytes, byte(orig&0xff))
		orig = orig >> 8
	}
	for len(bytes) != leng {
		bytes = append(bytes, byte(0x00))
	}
	return bytes
}
func IntToFloatStr(sats string, dec int) string {
	digits := len(sats)
	var ans string
	if sats == "0" {
		return "0.0"
	}
	if digits <= dec {
		zero_num := dec - digits
		zeros := strings.Repeat("0", zero_num)
		ans = "0." + zeros + sats
	} else {
		sats_byte := []byte(sats)
		big_byte := sats_byte[:len(sats_byte)-dec]
		small_byte := sats_byte[len(sats_byte)-dec:]
		big := string(big_byte)
		small := string(small_byte)
		ans = big + "." + small
	}
	ans_byte := []byte(ans)
	prev_length := len(ans_byte)
	for i := prev_length - 1; ans_byte[i] == byte('0'); i-- {
		ans_byte = ans_byte[:len(ans_byte)-1]
	}
	if ans_byte[len(ans_byte)-1] == byte('.') {
		ans_byte = append(ans_byte, byte('0'))
	}
	ans = string(ans_byte)
	return ans
}
func FloatStrToInt(btc string, dec int) int {
	sats := 0
	scale := 1
	for range dec {
		scale *= 10
	}
	if strings.Contains(btc, ".") {
		numl := strings.Split(btc, ".")
		big, _ := strconv.Atoi(numl[0])
		small_str := numl[1] + strings.Repeat("0", dec-len(numl[1]))
		small, _ := strconv.Atoi(small_str)
		sats += big * scale
		sats += small
	} else {
		am, _ := strconv.Atoi(btc)
		sats += am * scale
	}
	return sats
}
func LampsToSol(sats string) string {
	return IntToFloatStr(sats, 9)
}

func SolToLamps(btc string) int {
	return FloatStrToInt(btc, 9)
}

func SaveTx(tx Tx) {
	MkdirOrNothing("temp")
	f, _ := os.OpenFile("temp/SOL_transaction.json", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	save, _ := json.MarshalIndent(tx, "", "    ")
	f.Write(save)
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func StrToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
