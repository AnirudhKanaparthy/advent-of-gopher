package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
)

const bytesInMD5Hash int = 16
const coinSecretKey string = "yzbqklnj"
const miningProgressFrequency int = 100000

type AdventCoin struct {
	secretKey string
	nounce    int
}

type AdventCoinError struct {
	s string
}

func (ac *AdventCoinError) Error() string {
	return ac.s
}

func MakeAdventCoin(secretKey string) AdventCoin {
	return AdventCoin{secretKey, 0}
}

func (ac *AdventCoin) Mine(difficulty int) error {
	// This is the classic proof of work introduced by HashCash, Bitcoin uses SHA256 instead of MD5 though

	if difficulty > bytesInMD5Hash {
		return &AdventCoinError{"difficulty is more than maximum number of bytes in MD5 hash"}
	}
	hashesChecked := 0
	ac.nounce = 0
	for {
		digest := md5.Sum([]byte(ac.secretKey + strconv.Itoa(ac.nounce)))
		hexDigest := hex.EncodeToString(digest[:])
		mined := true
		for i := 0; i < difficulty; i += 1 {
			if hexDigest[i] != '0' {
				mined = false
				break
			}
		}
		if mined {
			return nil
		}

		if hashesChecked%miningProgressFrequency == 0 {
			fmt.Printf("Number of hashes checked: %v, last hash: %v\n", hashesChecked, hexDigest)
		}
		hashesChecked += 1
		ac.nounce += 1
	}
}

func (coin *AdventCoin) PrintDetails() {
	data := coin.secretKey + strconv.Itoa(coin.nounce)
	digest := md5.Sum([]byte(data))

	fmt.Printf("{secretKey: %v, nounce: %v}\n", coin.secretKey, strconv.Itoa(coin.nounce))
	fmt.Printf("hash: %v\n", hex.EncodeToString(digest[:]))
	fmt.Printf("nounce: %v\n", coin.nounce)
}

func main() {
	coin := MakeAdventCoin(coinSecretKey)
	coin.Mine(5)
	coin.PrintDetails()

	fmt.Printf("\n")

	coin.Mine(6)
	coin.PrintDetails()
}
