package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

// Modified version of https://gist.github.com/josephspurrier/12cc5ed76d2228a41ceb

func Decrypt(ciphertext []byte, keyString string) ([]byte, error) {

	// Key
	hash := sha256.Sum256([]byte(keyString))
	key := hash[:]

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		return []byte(""), errors.New("ciphertext too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func Encrypt(plaintext []byte, keyString string) ([]byte, error) {

	// Key
	hash := sha256.Sum256([]byte(keyString))
	key := hash[:] // Convert [32]byte into []byte

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), err
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte(""), err
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func ConvertBytesToHex(bytes []byte) string {
	dst := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(dst, bytes)
	return string(dst)
}
