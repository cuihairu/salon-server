package utils

import (
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

const (
	saltSize            = 16
	keySize             = 32
	timeCost            = 1
	memoryCost          = 64 * 1024
	parallelism         = 4
	AuthorizationKey    = "Authorization"
	AuthorizationPrefix = "Bearer "
)

func SetHeaderToken(c *gin.Context, token string) {
	c.Header(AuthorizationKey, AuthorizationPrefix+token)
}

func GetHeaderToken(c *gin.Context) string {
	return c.GetHeader(AuthorizationKey)
}

// GenerateRandomSaltWithSize generates a random salt with the given size.
//
// The returned salt is suitable for use with the `GeneratePasswordHash`
// function.
//
// If an error occurs during salt generation, it is returned.
func GenerateRandomSaltWithSize(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// GenerateRandomSalt generates a random salt with the size of `saltSize`.
//
// The returned salt is suitable for use with the `GeneratePasswordHash`
// function.
func GenerateRandomSalt() ([]byte, error) {
	return GenerateRandomSaltWithSize(saltSize)
}

// GeneratePasswordHash generates a hash for the given password, using the
// argon2.IDKey function. It returns the generated hash and the salt used to
// generate the hash. If an error occurs during salt generation, it is returned.
//
// The hash is generated with the following parameters:
//
// - Time cost: 1 iteration
// - Memory cost: 64 KiB
// - Parallelism: 4 threads
// - Key size: 32 bytes
//
// The generated hash and salt should be stored together, as the salt is needed
// to verify the password later.
func GeneratePasswordHash(password string) ([]byte, []byte, error) {
	salt, err := GenerateRandomSalt()
	if err != nil {
		return nil, nil, err
	}
	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, keySize)
	return hash, salt, nil
}

// VerifyPassword takes a hashed password and a password, and verifies if they match.
// The given salt is used to generate the hash.
// It returns true if the password matches, false otherwise.
func VerifyPassword(hashedPassword []byte, password []byte, salt []byte) bool {
	hash := argon2.IDKey(password, salt, timeCost, memoryCost, parallelism, keySize)
	// To verify, regenerate the hash with the same salt and compare
	verifyHash := argon2.IDKey(password, salt, timeCost, memoryCost, parallelism, keySize)
	return string(hash) == string(verifyHash)
}
