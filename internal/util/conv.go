package util

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
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

func IntToBytes(value int, length int) ([]byte, error) {
	if length == 4 {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(value))
		return buf, nil
	} else if length == 8 {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(value))
		return buf, nil
	} else {
		return nil, errors.New("Length must be 4 or 8.")
	}
}

func LampsToSol(sats string) string {
	digits := len(sats)
	var ans string
	if sats == "0" {
		return "0.0"
	}
	if digits <= 9 {
		zero_num := 9 - digits
		zeros := strings.Repeat("0", zero_num)
		ans = "0." + zeros + sats
	} else {
		sats_byte := []byte(sats)
		big_byte := sats_byte[:len(sats_byte)-9]
		small_byte := sats_byte[len(sats_byte)-9:]
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

func SolToLamps(btc string) int {
	sats := 0
	if strings.Contains(btc, ".") {
		numl := strings.Split(btc, ".")
		big, _ := strconv.Atoi(numl[0])
		small_str := numl[1] + strings.Repeat("0", 9-len(numl[1]))
		small, _ := strconv.Atoi(small_str)
		sats += big * 1000000000
		sats += small
	} else {
		am, _ := strconv.Atoi(btc)
		sats += am * 1000000000
	}
	return sats
}

func SaveTx(tx Tx) {
	MkdirOrNothing("temp")
	f, _ := os.OpenFile("temp/transaction.json", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
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
