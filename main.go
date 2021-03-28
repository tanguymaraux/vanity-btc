package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"vanity-btc/cmd"

	"github.com/btcsuite/btcutil/base58"

	"go.dedis.ch/kyber/v3/suites"
)

func main() {
	cmd.Execute()
}

func verifyaddressCheckSum(bitcoin_address string) bool {
	bitcoinaddressByte := base58.Decode(bitcoin_address)
	bitcoinWithoutChecksum := bitcoinaddressByte[:21]

	firstHash := sha256.Sum256(bitcoinWithoutChecksum)
	secondHash := sha256.Sum256(firstHash[:])
	addresschecksum := secondHash[:4]

	return bytes.Equal(addresschecksum, bitcoinaddressByte[21:25])
}

func test() {
	s := suites.MustFind("Ed25519")
	x := s.Scalar().SetInt64(42)
	P := s.Point().Mul(x, s.Point().Base())
	fmt.Println(P.String())
}
