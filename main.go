package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v3"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	privatekey, _ := secp256k1.GeneratePrivateKey()
	fmt.Println(computeAdress(privatekey))
	fmt.Println(privatekey.Key.String())

	s := suites.MustFind("Ed25519")
	x := s.Scalar().SetInt64(42)
	P := s.Point().Mul(x, s.Point().Base())
	fmt.Println(P.String())
}

func computeAdress(privatekey *secp256k1.PrivateKey) string {
	publickey := privatekey.PubKey().SerializeCompressed()
	hash := sha256.Sum256(publickey)
	hash_ripemd160 := ripemd160.New()
	hash_ripemd160.Write(hash[:])
	var version_hash []byte = []byte{00}
	version_hash = append(version_hash, hash_ripemd160.Sum(nil)...)
	newhash := sha256.Sum256(version_hash)
	newhash2 := sha256.Sum256(newhash[:])
	adresschecksum := newhash2[:4]
	bitcoin_adress := append(version_hash, adresschecksum...)
	bitcoin_byte_format := base58.Encode(bitcoin_adress)

	return bitcoin_byte_format
}

func verifyAdress(bitcoin_adress string) bool {
	bitcoinAdressByte := base58.Decode(bitcoin_adress)
	bitcoinWithoutChecksum := bitcoinAdressByte[:21]
	firstHash := sha256.Sum256(bitcoinWithoutChecksum)
	secondHash := sha256.Sum256(firstHash[:])
	adresschecksum := secondHash[:4]

	return bytes.Equal(adresschecksum, bitcoinAdressByte[21:25])
}
