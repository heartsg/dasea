package common

import (
	"crypto/rand"
	"crypto/md5"
	"golang.org/x/crypto/scrypt"
)
const (
	saltSize = 16
	passwordSize = 32
	//default md5 size for golang
	md5Size = md5.Size
)

func Salt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

//According to benchmarking, scrypt.Key takes 60ms for i7 processor.
//Use encrypt only when necessary 
func Encrypt(password string, salt []byte) ([]byte, error) {
	key, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, passwordSize)
	if err != nil {
		return nil, err
	}
	return key, nil
}


func Md5Encrypt(password string) []byte {
	key := make([]byte, md5Size)
	staticKey := md5.Sum([]byte(password))
	for i := 0; i < md5Size; i++ {
		key[i] = staticKey[i]
	}
	return key
}