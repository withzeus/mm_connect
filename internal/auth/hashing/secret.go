package hashing

import "golang.org/x/crypto/bcrypt"

func Hash(secret string) (string, error) {
	b, err := bcrypt.GenerateFromPassword(
		[]byte(secret),
		bcrypt.DefaultCost,
	)

	return string(b), err
}

func Verify(hash, secret string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(secret),
	)
}
