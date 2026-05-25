package uuid

import "github.com/google/uuid"

func GenerateV4() string {
	return uuid.New().String()
}
