package main

import (
	"encoding/hex"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v3"
	"github.com/stretchr/testify/assert"
)

var privatekey = "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"

func TestComputeAdress(t *testing.T) {
	key, err := hex.DecodeString(privatekey)
	assert.Nil(t, err)
	privatekey := secp256k1.PrivKeyFromBytes(key)

	assert.Equal(t, "1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs", computeAdress(privatekey))
}

func TestVerifyChecksum(t *testing.T) {
	assert.True(t, verifyAdress("1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs"))
	assert.False(t, verifyAdress("1PMycacNJaSqwwJqjawXBErnLsZ7RkXUAs"))
}
