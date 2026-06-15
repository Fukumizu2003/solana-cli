package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/scrypt"
)

func B58Encode(msg []byte) string {
	return base58.Encode(msg)
}

func B58Decode(msg string) []byte {
	return base58.Decode(msg)
}

func GenKey(len int) []byte {
	buf := make([]byte, len)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}
func NewKeypair() (ed25519.PrivateKey, ed25519.PublicKey) {
	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	return privKey, pubKey
}

func BytesToKeypair(priv []byte) (ed25519.PrivateKey, ed25519.PublicKey) {
	privKey := ed25519.NewKeyFromSeed(priv)
	pubKey := privKey.Public().(ed25519.PublicKey)
	return privKey, pubKey
}

func PubkeyToAddress(pub []byte) string {
	address := B58Encode(pub)
	return address
}

func AddressToPubkey(addr string) []byte {
	pubk := B58Decode(addr)
	return pubk
}

func ScryptHashNew(pw []byte) ([]byte, []byte) {
	salt := GenKey(16) // Always use a unique salt
	N := 16384         // CPU/memory cost parameter
	r := 8             // Block size
	p := 1             // Parallelization factor
	keyLen := 32
	key, err := scrypt.Key(pw, salt, N, r, p, keyLen)
	if err != nil {
		log.Fatalf("Error generating key: %v", err)
	}
	return salt, key
}

func ScryptHashAgain(salt []byte, pw []byte) []byte {
	N := 16384 // CPU/memory cost parameter
	r := 8     // Block size
	p := 1     // Parallelization factor
	keyLen := 32
	key, err := scrypt.Key(pw, salt, N, r, p, keyLen)
	if err != nil {
		log.Fatalf("Error generating key: %v", err)
	}
	return key
}

func AesEncrypt(priv []byte, pw []byte) []byte {
	salt, key := ScryptHashNew(pw)
	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)

	nonce := GenKey(12)
	rand.Read(nonce)
	ciphertext := aesgcm.Seal(nil, nonce, priv, nil)
	ans := append(append(salt, nonce...), ciphertext...)
	return ans
}

func AesDecrypt(priv []byte, pw []byte) ([]byte, error) {
	salt := priv[:16]
	nonce := priv[16:28]
	key := ScryptHashAgain(salt, pw)
	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)

	decrypted, err := aesgcm.Open(nil, nonce, priv[28:], nil)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
