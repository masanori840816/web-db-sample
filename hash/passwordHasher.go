package hash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"

	"math/big"
)

// itelateCount(4) + salt length(4)
const fixedPasswordLength = 8

func GeneratePasswordHash(password string) (string, error) {
	salt, err := generateRandomSalt(128 / 8)
	if err != nil {
		return "", err
	}
	result := generateHash(password, salt, 100_000, 256/8)
	return base64.URLEncoding.EncodeToString(result), nil
}
func VerifyPassword(inputPassword string, hashedPassword string) (bool, error) {
	decodedPassword, err := base64.URLEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, err
	}
	iterateCount := readNetworkByteOrder(decodedPassword, 0)
	saltLength := readNetworkByteOrder(decodedPassword, 4)
	salt := make([]byte, saltLength)
	blockCopy(decodedPassword, fixedPasswordLength, salt, 0, int(saltLength))
	hashedInput := generateHash(inputPassword, salt, int(iterateCount),
		len(decodedPassword)-(int(saltLength)+fixedPasswordLength))
	return bytes.Equal(hashedInput, decodedPassword), nil
}
func generateHash(original string, salt []byte, iterateCount int, keyLength int) []byte {
	// Get base 64 encoded Hasu value to save the password
	key := pbkdf2.Key([]byte(original), salt, iterateCount, keyLength, sha512.New)
	// Add the iterate count, salt, salt length.
	results := make([]byte, len(key)+fixedPasswordLength+len(salt))
	writeNetworkByteOrder(results, 0, uint(iterateCount))
	writeNetworkByteOrder(results, 4, uint(len(salt)))
	// Add salt
	blockCopy(salt, 0, results, fixedPasswordLength, len(salt))
	// Add hashed password
	blockCopy(key, 0, results, fixedPasswordLength+len(salt), len(key))
	return results
}

// Generate a salt value
func generateRandomSalt(length int) ([]byte, error) {
	results := make([]byte, length)
	for i := 0; i < length; i++ {
		salt, err := rand.Int(rand.Reader, big.NewInt(255))
		if err != nil {
			return nil, err
		}
		results[i] = byte(salt.Int64())
	}
	return results, nil
}
func writeNetworkByteOrder(buffer []byte, offset int, value uint) {
	buffer[offset+0] = byte(value >> 24)
	buffer[offset+1] = byte(value >> 16)
	buffer[offset+2] = byte(value >> 8)
	buffer[offset+3] = byte(value >> 0)
}
func readNetworkByteOrder(buffer []byte, offset int) uint {
	return ((uint)(buffer[offset]) << 24) |
		((uint)(buffer[offset+1]) << 16) |
		((uint)(buffer[offset+2]) << 8) |
		((uint)(buffer[offset+3]))
}
func blockCopy(src []byte, srcOffset int, dst []byte, dstOffset int, copyLength int) {
	index := dstOffset
	for i := srcOffset; i < copyLength+srcOffset; i++ {
		dst[index] = src[i]
		index += 1
	}
}
