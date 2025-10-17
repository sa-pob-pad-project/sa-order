package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"


	"golang.org/x/crypto/argon2"
)

const (
	Memory      = 64 * 1024
	Iterations  = 3
	Parallelism = 4
	KeyLength   = 32
)

func GenerateRandomByte(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HashPassword(password string) (string, error) {
	salt, err := GenerateRandomByte(16)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, Iterations, Memory, Parallelism, KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, Memory, Iterations, Parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	dbPasswordHash := strings.Split(encodedHash, "$")[5]
	vals := strings.Split(encodedHash, "$")
	decodedSalt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false, err
	}
	userInputPaswordHash := argon2.IDKey([]byte(password), decodedSalt, Iterations, Memory, Parallelism, KeyLength)
	base64InputPasswordHash := base64.RawStdEncoding.EncodeToString(userInputPaswordHash)
	return dbPasswordHash == base64InputPasswordHash, nil
}
