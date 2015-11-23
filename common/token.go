package common

import (
	"crypto/rand"
)

func Token(size int) (string, error) {
	token := make([]byte, size)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for k, v := range token {
		token[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(token), nil
}