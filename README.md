## Table of contents
* [Vanity-btc](#vanity-btc)
* [Technologies](#technologies)
* [Usage](#usage)
* [Generation Process](#generation-process)
* [Example](#example)
* [Docker](#docker)
* [License](#license)

## Vanity-btc
Vanity-btc allows you to create vanity bitcoin addresses using golang.
	
## Technologies
This project is created with:
* [Golang](https://golang.org/) version 1.16.2

Using packages:
* github.com/spf13/cobra
* go.uber.org/zap
* github.com/decred/base58
* [github.com/decred/dcrd/dcrec/secp256k1/v3](https://pkg.go.dev/github.com/decred/dcrd/dcrec/secp256k1/v3)
* golang.org/x/crypto/ripemd160
	
## Usage
```
$ go build
$ ./vanity-btc --help
vanity-btc is a bitcoin vanity address generator.

It is based on this method : https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addressses
Vanity-btc version 1.0

Usage:
  vanity-btc [flags]

Flags:
  -c, --chronometer      enable chronometer
  -h, --help             help for vanity-btc
  -n, --number           enable counting the number of addresses generated
  -p, --pattern string   pattern wanted in the btc address
  -t, --threads int      number of threads to use (default 4)
  -v, --verbose          enable verbose mode
```

## Generation Process
The generation method is based on [this principle](https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses)

**secp256k1** → **ECDSA Private Key** → public key<br>
public key → **crypto/sha256** → **ripemd160**<br>
Add **version** byte (0x00 for Main Network)<br>
**extended RIPEMD-160** → **crypto/sha256** → **crypto/sha256**<br>
Take the first 4 bytes<br>
Add the 4 bytes at the end of extended RIPEMD-160 hash<br>
**base58** → bitcoin address<br>

## Example
```
$ ./vanity-btc --chronometer --number --pattern "BTC" --threads 12
2021-03-29T10:50:21.729+0100	INFO	addresses/pattern.go:54	Generating address with patern: "BTC"...
2021-03-29T10:50:22.259+0100	INFO	addresses/pattern.go:66	Address found!
2021-03-29T10:50:22.259+0100	INFO	addresses/pattern.go:70	Number of address generated: 45167
2021-03-29T10:50:22.259+0100	INFO	addresses/pattern.go:74	
Address: 1BTC4MjELFrYS8mCvnJ3RHgqBt4J9L8y84
Private Key: e5246a1d3631df439979d94ad7dd2bfd4998efc1b777b1de4b5d0959b88afb7a
Execution: 528.387659ms
```

## Docker
```
docker run -it --rm pilpur/vanity-btc --chronometer --number --pattern "BTC" --threads 12
```

## License
This software is licensed under [MIT license.](LICENSE)<br>
Copyright (c) 2021 Tanguy Maraux