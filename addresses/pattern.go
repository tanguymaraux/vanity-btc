package addresses

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/decred/base58"
	"github.com/decred/dcrd/dcrec/secp256k1/v3"
	"golang.org/x/crypto/ripemd160"
)

var (
	stop  bool
	done      = make(chan result)
	count int = 0
)

type result struct {
	address    string
	privatekey string
}

func GetPattern(pattern string, verbose bool, threads int, difficulty bool) {
	if !verbose {
		fmt.Println("Generating address...\n")
	}

	for i := 0; i < threads; i++ {
		go generateaddressWithPattern(pattern, verbose, difficulty)
	}

	r := <-done

	fmt.Println("Address found!")
	if difficulty {
		fmt.Println("Difficulty:", count)
	}
	fmt.Println("Address:", r.address, "\nPrivate Key:", r.privatekey)

}

func CheckPattern(pattern string) bool {
	var alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	if len(pattern) > 33 {
		return false
	}

	for _, char := range pattern {
		if !(strings.ContainsRune(alphabet, char)) {
			return false
		}
	}

	return true
}

func generateaddressWithPattern(pattern string, verbose bool, difficulty bool) {
	address, privkey := "", ""
	for {
		if stop {
			return
		}
		address, privkey = Generateaddress()
		if difficulty {
			count++
		}
		if strings.HasPrefix(address[1:], pattern) {
			break
		}
		if verbose {
			fmt.Println(address)
		}
	}

	done <- result{address: address, privatekey: privkey}
	stop = true
}

func Generateaddress() (string, string) {
	privatekey, _ := secp256k1.GeneratePrivateKey()
	return computeaddress(privatekey), privatekey.Key.String()
}

func computeaddress(privatekey *secp256k1.PrivateKey) string {
	publickey := privatekey.PubKey().SerializeCompressed()
	hash := sha256.Sum256(publickey)

	hash_ripemd160 := ripemd160.New()
	hash_ripemd160.Write(hash[:])

	var version_hash []byte = []byte{00}
	version_hash = append(version_hash, hash_ripemd160.Sum(nil)...)

	newhash := sha256.Sum256(version_hash)
	newhash2 := sha256.Sum256(newhash[:])
	addresschecksum := newhash2[:4]

	bitcoin_address := append(version_hash, addresschecksum...)

	bitcoin_byte_format := base58.Encode(bitcoin_address)

	return bitcoin_byte_format
}
