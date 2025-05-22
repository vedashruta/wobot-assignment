package bcrypt

import "golang.org/x/crypto/bcrypt"

func Hash(pwd string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	if err != nil {
		return
	}
	hash = string(bytes)
	return
}

func Verify(password string, hash string) (ok bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return
	}
	ok = true
	return
}
