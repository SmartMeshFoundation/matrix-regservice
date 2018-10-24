package utils

import (
	"bytes"
	"crypto/ecdsa"
	rand2 "crypto/rand"
	"encoding/base32"
	"encoding/gob"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//EmptyHash all zero,invalid
var EmptyHash = common.Hash{}

//EmptyAddress all zero,invalid
var EmptyAddress = common.Address{}

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

//Pex short string stands for data
func Pex(data []byte) string {
	return common.Bytes2Hex(data[:4])
}

//HPex pex for hash
func HPex(data common.Hash) string {
	return common.Bytes2Hex(data[:2])
}

//APex pex for address
func APex(data common.Address) string {
	return common.Bytes2Hex(data[:4])
}

//APex2 shorter than APex
func APex2(data common.Address) string {
	return common.Bytes2Hex(data[:2])
}

//PubkeyToAddress convert pubkey bin to address
func PubkeyToAddress(pubkey []byte) common.Address {
	return common.BytesToAddress(crypto.Keccak256(pubkey[1:])[12:])
}

//BigIntTo32Bytes convert a big int to bytes
func BigIntTo32Bytes(i *big.Int) []byte {
	data := i.Bytes()
	buf := make([]byte, 32)
	for i := 0; i < 32-len(data); i++ {
		buf[i] = 0
	}
	for i := 32 - len(data); i < 32; i++ {
		buf[i] = data[i-32+len(data)]
	}
	return buf
}

// BytesToString accepts bytes and returns their string presentation
// instead of string() this method doesn't generate memory allocations,
// BUT it is not safe to use anywhere because it points
// this helps on 0 memory allocations
func BytesToString(b []byte) string {
	/* #nosec */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: bh.Data,
		Len:  bh.Len,
	}
	/* #nosec */
	return *(*string)(unsafe.Pointer(&sh))
}

// StringToBytes accepts string and returns their []byte presentation
// instead of byte() this method doesn't generate memory allocations,
// BUT it is not safe to use anywhere because it points
// this helps on 0 memory allocations
func StringToBytes(s string) []byte {
	/* #nosec */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  0,
	}
	/* #nosec */
	return *(*[]byte)(unsafe.Pointer(&bh))
}

//
const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

//RandSrc random source from math
var RandSrc = rand.NewSource(time.Now().UnixNano())

func readFullOrPanic(r io.Reader, v []byte) int {
	n, err := io.ReadFull(r, v)
	if err != nil {
		panic(err)
	}
	return n
}

// Random takes a parameter (int) and returns random slice of byte
// ex: var randomstrbytes []byte; randomstrbytes = utils.Random(32)
func Random(n int) []byte {
	v := make([]byte, n)
	readFullOrPanic(rand2.Reader, v)
	return v
}

// RandomString accepts a number(10 for example) and returns a random string using simple but fairly safe random algorithm
func RandomString(n int) string {
	s := base32.StdEncoding.EncodeToString(Random(n))
	return s[:n]
}

//NewRandomInt generate a random int ,not more than n
func NewRandomInt(n int) int {
	return rand.New(RandSrc).Intn(n)
}

//NewRandomInt64 generate a random int64
func NewRandomInt64() int64 {
	return rand.New(RandSrc).Int63()
}

//NewRandomAddress generate a address,there maybe no corresponding priv key
func NewRandomAddress() common.Address {
	hash := Sha3([]byte(Random(10)))
	return common.BytesToAddress(hash[12:])
}

//MakePrivateKeyAddress generate a private key and it's address
func MakePrivateKeyAddress() (*ecdsa.PrivateKey, common.Address) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}

//StringInterface use spew to string any object with max `depth`,it's not thread safe.
func StringInterface(i interface{}, depth int) string {
	stringer, ok := i.(fmt.Stringer)
	if ok {
		return stringer.String()
	}
	c := spew.Config
	spew.Config.DisableMethods = false
	//spew.Config.ContinueOnMethod = false
	spew.Config.MaxDepth = depth
	s := spew.Sdump(i)
	spew.Config = c
	return s
}

//StringInterface1 use spew to string any object with depth 1,it's not thread safe.
func StringInterface1(i interface{}) string {
	stringer, ok := i.(fmt.Stringer)
	if ok {
		return stringer.String()
	}
	c := spew.Config
	spew.Config.DisableMethods = true
	spew.Config.MaxDepth = 1
	s := spew.Sdump(i)
	spew.Config = c
	return s
}

//DeepCopy use gob to copy
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// Exists returns true if directory||file exists
func Exists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// SecretGenerator Return n random bytes suitable for cryptographic and it's corresponding hashlock
type SecretGenerator func() (secret, hashlock common.Hash)

//RandomSecretGenerator generate a random secret
func RandomSecretGenerator() (secret, hashlock common.Hash) {
	secret = common.BytesToHash(Random(32))
	hashlock = Sha3(secret[:])
	return
}

//NewSepecifiedSecretGenerator create a secret generator which don't know the secret
func NewSepecifiedSecretGenerator(hashlock common.Hash) SecretGenerator {
	return func() (secret, hashlock2 common.Hash) {
		return EmptyHash, hashlock
	}
}

// GetHomePath returns the user's $HOME directory
func GetHomePath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	return os.Getenv("HOME")
}

//SystemExit quit and print stack
func SystemExit(code int) {
	if code != 0 {
		debug.PrintStack()
	}
	time.Sleep(time.Second * 2)
	os.Exit(code)
}
