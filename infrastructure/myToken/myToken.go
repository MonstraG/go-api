package myToken

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/models"
	"io"
	"time"
	"uuid"
)

const Cookie = "token"
const DefaultCookieAge = 3600 * 24
const defaultTokenExpiry = DefaultCookieAge * time.Second

type Service struct {
	// cipherBlock is the aes.NewCipher, re-used between calls
	cipherBlock cipher.Block

	// now is the source for time
	now func() time.Time

	// ivReader is the source for iv in AES
	ivReader io.Reader

	// maxAge is after what time does the token expire
	maxAge time.Duration
}

func NewService(config appConfig.AppConfig) (Service, error) {
	secretBytes, err := hex.DecodeString(config.TokenSecret)
	if err != nil {
		return Service{}, fmt.Errorf("failed to decode token secret: %v", err)
	}

	cipherBlock, err := aes.NewCipher(secretBytes)
	if err != nil {
		return Service{}, err
	}

	return Service{
		now:         time.Now,
		cipherBlock: cipherBlock,
		ivReader:    rand.Reader,
		maxAge:      defaultTokenExpiry,
	}, nil
}

type TokenPayload struct {
	Sub uuid.UUID `json:"sub"`
	Iat int64     `json:"iat"`
}

func (myToken *Service) CreateToken(user models.User) (string, error) {
	payload := TokenPayload{
		Sub: user.ID,
		Iat: myToken.now().Unix(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	ciphertext, err := myToken.aes256CbcEncode(payloadBytes)
	if err != nil {
		return "", err
	}

	// and now base64 just for easy copy-paste
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (myToken *Service) aes256CbcEncode(plaintext []byte) ([]byte, error) {
	blockSize := myToken.cipherBlock.BlockSize()

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block.
	paddedText := pkcs7Pad(plaintext, blockSize)

	// The IV needs to be unique, but not secure. Therefore, it's common to
	// include it at the beginning of the ciphertext. Therefore, we add another `blockSize` for that IV.
	output := make([]byte, blockSize+len(paddedText))
	iv := output[:blockSize]
	ciphertext := output[blockSize:]

	if _, err := io.ReadFull(myToken.ivReader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(myToken.cipherBlock, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	return output, nil
}

func pkcs7Pad(plaintext []byte, blockSize int) []byte {
	paddingSize := blockSize - len(plaintext)%blockSize
	if paddingSize == 0 {
		return plaintext
	}

	paddedText := make([]byte, len(plaintext)+paddingSize)
	copy(paddedText, plaintext)
	copy(paddedText[len(plaintext):], bytes.Repeat([]byte{byte(paddingSize)}, paddingSize))

	return paddedText
}

func pkcs7Unpad(paddedText []byte, blockSize int) ([]byte, error) {
	pad := paddedText[len(paddedText)-1]
	if pad < 1 || pad > byte(blockSize) {
		return nil, ErrPadInvalid
	}

	padStart := len(paddedText) - int(pad)

	for _, padByte := range paddedText[padStart:] {
		if padByte != pad {
			return nil, ErrPadInvalid
		}
	}

	return paddedText[:padStart], nil
}

func (myToken *Service) aes256CbcDecode(ciphertext []byte) ([]byte, error) {
	blockSize := myToken.cipherBlock.BlockSize()

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(myToken.cipherBlock, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	ciphertext, err := pkcs7Unpad(ciphertext, blockSize)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func (myToken *Service) ParseToken(tokenString string) (TokenPayload, error) {
	tokenPayload := TokenPayload{}

	tokenBytes, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return tokenPayload, fmt.Errorf("token base64 decode failed: %v", err)
	}

	plaintext, err := myToken.aes256CbcDecode(tokenBytes)
	if err != nil {
		return tokenPayload, fmt.Errorf("token aes decode failed: %v", err)
	}

	err = json.Unmarshal(plaintext, &tokenPayload)
	if err != nil {
		return tokenPayload, fmt.Errorf("failed unmarshal: %v", err)
	}

	if tokenPayload.Sub == uuid.Nil() {
		return tokenPayload, ErrTokenInvalid
	}

	expiryAt := time.Unix(tokenPayload.Iat, 0).Add(myToken.maxAge)
	if time.Now().After(expiryAt) {
		return tokenPayload, ErrTokenExpired
	}

	return tokenPayload, nil
}

var (
	ErrPadInvalid   = errors.New("padding is invalid")
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
)
