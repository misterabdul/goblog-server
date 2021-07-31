package hash

import "golang.org/x/crypto/bcrypt"

// Create bcrypt's hashed string from given password string.
func Make(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

// Check the given password against bcrypt's hashed string.
func Check(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
