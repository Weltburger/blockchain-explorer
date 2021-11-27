package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func encryptPwd(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	encryptedPwd, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPwd := fmt.Sprintf("%s.%s", hex.EncodeToString(encryptedPwd), hex.EncodeToString(salt))

	return hashedPwd, nil
}

func comparePwd(storedPwd string, suppliedPwd string) (bool, error) {
	pwSalt := strings.Split(storedPwd, ".")

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwSalt[1])
	if err != nil {
		return false, err
	}

	encryptedPwd, err := scrypt.Key([]byte(suppliedPwd), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, err
	}

	return hex.EncodeToString(encryptedPwd) == pwSalt[0], nil
}
