package hash

import "golang.org/x/crypto/bcrypt"

// Create bcrypt's hashed string from given password string.
func Make(password string) (hashedPassword string, err error) {
	var hashedPasswordByte []byte

	if hashedPasswordByte, err = bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	); err != nil {
		return "", err
	}

	return string(hashedPasswordByte), nil
}

// Check the given password against bcrypt's hashed string.
func Check(password string, hashedPassword string) (match bool) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password))

	return err == nil
}
