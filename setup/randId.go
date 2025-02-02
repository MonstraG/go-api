package setup

import (
	"log"
	"math/rand"
	"strings"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

const (
	// 6 bits are enough to represent 64 indexes (6 = 111111 in binary)
	bitsForAlphabet = 6
	// the bits 0..0111111, exactly the bitsForAlphabet value in binary
	bitMask = 1<<bitsForAlphabet - 1
	// number of letter indices fitting in 63 bits
	lettersPerRandomCall = 63 / bitsForAlphabet
	// this can be increased up to lettersPerRandomCall
	lettersToMake = 6
)

var stringBuilder = strings.Builder{}

// This function checks that const make sense and will not break the id generation.
// Inspection is disabled so that it wouldn't tell me that everything is fine.
//
//goland:noinspection GoBoolExpressions
func init() {
	stringBuilder.Grow(lettersToMake)

	if lettersToMake > lettersPerRandomCall {
		log.Fatalf("You cannot generate so many letters from one int63 call, max is %d", lettersPerRandomCall)
	}
	if len(alphabet) < bitMask {
		log.Fatalf("Alphabet is too small, and must have at least %d letters", bitMask)
	}
}

// RandId generates a random id of length lettersToMake
//
// Based on https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go,
// but with constants adjusted so that safety checks can be removed (e.g. alphabet large enough to avoid out-of bounds)
func RandId() string {
	stringBuilder.Reset()
	randomBits := rand.Int63()

	for lettersRemain := lettersToMake; lettersRemain > 0; {
		// take part of the bits and figure out letter index
		letterIndex := int(randomBits & bitMask)
		// get rid of these bits by shifting everything to the right
		randomBits >>= bitsForAlphabet
		// add this letter
		stringBuilder.WriteByte(alphabet[letterIndex])
		lettersRemain--
	}

	return stringBuilder.String()
}
