package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password, passwordPepper string, saltRounds int) []byte {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), saltRounds)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	fmt.Println(err) // nil means it is a match
	return hashedPassword
}

//HashPasswordWrapper calc password hash for matrix homeserver
func HashPasswordWrapper(password string) string {
	hash := hashPassword(password, "", 12)
	return string(hash)
}

//VerifyPassword test password and password_hash matches
func VerifyPassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
