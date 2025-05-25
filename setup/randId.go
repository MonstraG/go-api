package setup

import (
	"go-api/setup/myLog"
	"math/rand"
	"strings"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

const (
	// 6 bits are enough to represent 64 indexes (6 = 111111 in binary)
	bitsForLetter = 6
	// all bits set
	bitMask = 1<<bitsForLetter - 1
	// I use int63 to generate random bits, it has 63 bits
	randomCallBitSize = 63
	// number of letter indices fitting in 63 bits
	lettersPerRandomCall = randomCallBitSize / bitsForLetter
	// this can be increased up to lettersPerRandomCall
	idLength = 6
)

var stringBuilder = strings.Builder{}

// This function checks that const make sense and will not break the id generation.
// Inspection is disabled so that it wouldn't tell me that everything is fine.
//
//goland:noinspection GoBoolExpressions
func init() {
	stringBuilder.Grow(idLength)

	if idLength > lettersPerRandomCall {
		myLog.Fatal.Logf("You cannot generate so many letters from one int63 call, max is %d", lettersPerRandomCall)
	}
	if len(alphabet) < bitMask {
		myLog.Fatal.Logf("Alphabet is too small, and must have at least %d letters", bitMask+1)
	}
}

// RandId generates a random id of length lettersToMake
//
// Based on https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go,
// but with constants adjusted so that safety checks can be removed (e.g. alphabet large enough to avoid out-of bounds)
func RandId() string {
	stringBuilder.Reset()
	randomBits := rand.Int63()

	for lettersRemain := idLength; lettersRemain > 0; {
		// take part of the bits and figure out letter index
		letterIndex := int(randomBits & bitMask)
		// get rid of these bits by shifting everything to the right
		randomBits >>= bitsForLetter
		// add this letter
		stringBuilder.WriteByte(alphabet[letterIndex])
		lettersRemain--
	}

	return stringBuilder.String()
}
