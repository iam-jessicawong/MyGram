package helpers

import "golang.org/x/crypto/bcrypt"

// password hashing before save to database
func HashPass(p string) string {
	salt := 8
	password := []byte(p)
	hash, _ := bcrypt.GenerateFromPassword(password, salt)

	return string(hash)
}

// compare password from database with user input from login
func ComparePass(h, p []byte) bool {
	hash, pass := []byte(h), []byte(p)
	err := bcrypt.CompareHashAndPassword(hash, pass)
	return err == nil
}
