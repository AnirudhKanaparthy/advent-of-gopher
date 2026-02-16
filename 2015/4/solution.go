package y2015d4

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const bytesInMD5Hash int = 16

type Solution struct{}

func MakeSolution() *Solution {
	return &Solution{}
}

func (sol Solution) ArgsString(int, []string) string {
	return "<secret key>"
}

func (sol Solution) Solve(part int, args []string, w io.Writer) error {
	if len(args) < 1 {
		return errors.New("No secret key provided")
	}

	coinSecretKey := args[0]

	coin := makeAdventCoin(coinSecretKey)
	switch part {
	case 1:
		coin.mine(5)
		coin.printDetails(w)
	case 2:
		coin.mine(6)
		coin.printDetails(w)
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	return nil
}

type adventCoin struct {
	secretKey string
	nounce    int
}

func makeAdventCoin(secretKey string) adventCoin {
	return adventCoin{secretKey, 0}
}

func (ac *adventCoin) mine(difficulty int) error {
	// This is the classic proof of work introduced by HashCash, Bitcoin uses SHA256 instead of MD5 though
	if difficulty > bytesInMD5Hash {
		return errors.New("difficulty is more than maximum number of bytes in MD5 hash")
	}
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
		ac.nounce += 1
	}
}

func (coin *adventCoin) printDetails(w io.Writer) {
	data := coin.secretKey + strconv.Itoa(coin.nounce)
	digest := md5.Sum([]byte(data))

	fmt.Fprintf(w, "secretKey : %v\n", coin.secretKey)
	fmt.Fprintf(w, "     hash : %v\n", hex.EncodeToString(digest[:]))
	fmt.Fprintf(w, "   nounce : %v", coin.nounce)
}
