package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password, passwordPepper string, saltRounds int) []byte {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), saltRounds)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	fmt.Println(err) // nil means it is a match
	return hashedPassword
}

//HashPasswordWrapper calc password hash for matrix homeserver
func HashPasswordWrapper(password string) string {
	hash := hashPassword(password, "", 12)
	return string(hash)
}

//VerifyPassword test password and password_hash matches
func VerifyPassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

//Ecrecover is a wrapper for crypto.Ecrecover
func Ecrecover(hash common.Hash, signature []byte) (addr common.Address, err error) {
	if len(signature) != 65 {
		err = fmt.Errorf("signature errr, len=%d,signature=%s", len(signature), hex.EncodeToString(signature))
		return
	}
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		signature[len(signature)-1] += 27
		return
	}
	addr = PubkeyToAddress(pubkey)
	signature[len(signature)-1] += 27
	return
}

//PubkeyToAddress convert pubkey bin to address
func PubkeyToAddress(pubkey []byte) common.Address {
	return common.BytesToAddress(crypto.Keccak256(pubkey[1:])[12:])
}

//SignData sign with ethereum format
func SignData(privKey *ecdsa.PrivateKey, data []byte) (sig []byte, err error) {
	hash := Sha3(data)
	//why add 27 for the last byte?
	sig, err = crypto.Sign(hash[:], privKey)
	if err == nil {
		sig[len(sig)-1] += byte(27)
	}
	return
}

//Sha3 is short for Keccak256Hash
func Sha3(data ...[]byte) common.Hash {
	return crypto.Keccak256Hash(data...)
}

//MakePrivateKeyAddress generate a private key and it's address
func MakePrivateKeyAddress() (*ecdsa.PrivateKey, common.Address) {
	//#nosec
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}
