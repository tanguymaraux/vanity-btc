// pattern.go - Main loop for generating address
// Copyright (c) 2021 Tanguy Maraux. Author Tanguy Maraux All rights reserved.
// Use of this source code is governed by a MIT license.

package addresses

import (
	"crypto/sha256"
	"strings"
	"time"

	"github.com/decred/base58"
	"github.com/decred/dcrd/dcrec/secp256k1/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// Init logger with debug and info mode, whe
func initlogger(verbose bool) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if verbose {
		// Debug mode
		config.Level.SetLevel(zap.DebugLevel)
	} else {
		// Info mode
		config.Level.SetLevel(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		zap.S().Panic(err)
	}
	zap.ReplaceGlobals(logger)
}

// Get the bitcoin vanity address
func GetAddress(pattern string, verbose bool, threads int, number bool) {
	// Init the logger
	initlogger(verbose)
	zap.S().Info("Generating address with patern: \"", pattern, "\"...")

	// Loop for multithreading the address's generating
	for i := 0; i < threads; i++ {
		// launch go routine
		go generateAddressWithPattern(pattern, verbose, number)
	}

	r := <-done

	// Wait for goroutines to end in order to not get parasite print in debug mode
	time.Sleep(time.Millisecond)
	zap.S().Info("Address found!")

	// Print the number of addresses generated
	if number {
		zap.S().Info("Number of address generated: ", count)
	}

	// Print the bitcoin vanity address and the private key related
	zap.S().Info("\nAddress: ", r.address, "\nPrivate Key: ", r.privatekey)

}

// Check if the pattern is valid
func CheckPattern(pattern string) bool {
	// Base58 alphabet
	var alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	// Check if every pattern's char is in the alphabet
	for _, char := range pattern {
		if !(strings.ContainsRune(alphabet, char)) {
			return false
		}
	}

	return true
}

// Generate an address and check if it contains the patern
func generateAddressWithPattern(pattern string, verbose bool, number bool) {
	address, privkey := "", ""
	for {
		// Stop if a goroutine found a valid address
		if stop {
			return
		}

		// Generate an address with its private key
		address, privkey = GenerateAddress()

		// Increment counter
		if number {
			count++
		}

		// Check if the pattern is at the beginning of the address
		if strings.HasPrefix(address[1:], pattern) {
			break
		}

		zap.S().Debug(address)
	}

	// Stores the values
	done <- result{address: address, privatekey: privkey}

	// Stops other goroutine
	stop = true
}

// Generate a random bitcoin address and its private key
func GenerateAddress() (string, string) {
	privatekey, _ := secp256k1.GeneratePrivateKey()

	// Compute the address from the private key
	return computeAddress(privatekey), privatekey.Key.String()
}

// Compute the bitcoin address
func computeAddress(privatekey *secp256k1.PrivateKey) string {
	// Get a private ECDSA key
	publickey := privatekey.PubKey().SerializeCompressed()

	// Perform SHA-256 hashing on the public key
	hash := sha256.Sum256(publickey)

	// Perform RIPEMD-160 hashing on the result of SHA-256
	hash_ripemd160 := ripemd160.New()
	hash_ripemd160.Write(hash[:])

	// Create the version byte (0x00 for Main Network)
	var version_hash []byte = []byte{00}
	// Add version byte in front of RIPEMD-160 hash
	version_hash = append(version_hash, hash_ripemd160.Sum(nil)...)

	// Get the address CheckSsum :
	// Perform SHA-256 hash on the extended RIPEMD-160 result
	newhash := sha256.Sum256(version_hash)
	// Perform SHA-256 hash on the result of the previous SHA-256 hash
	newhash2 := sha256.Sum256(newhash[:])
	// Take the address checksum (first 4 bytes of the second SHA-256 hash)
	addresschecksum := newhash2[:4]

	// Add the 4 checksum bytes at the end of extended RIPEMD-160 hash. This is the 25-byte binary Bitcoin Address.
	bitcoin_address := append(version_hash, addresschecksum...)
	// Convert the result from a byte string into a base58 string
	bitcoin_byte_format := base58.Encode(bitcoin_address)

	return bitcoin_byte_format
}
