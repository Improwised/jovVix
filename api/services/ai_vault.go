package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Improwised/jovvix/api/models"
	"golang.org/x/crypto/argon2"
)

var ErrAIVaultPassword = errors.New("the Vault Password is incorrect or saved data is invalid")

func vaultKey(p, s []byte) []byte { return argon2.IDKey(p, s, 3, 64*1024, 1, 32) }
func wipe(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
func ValidateAIVaultPassword(p string) error {
	if len(p) < 8 || len(p) > 512 {
		return fmt.Errorf("Vault Password must be between 8 and 512 characters")
	}
	return nil
}
func EncryptAIKey(user string, password, key []byte) ([]byte, []byte, []byte, error) {
	defer wipe(password)
	defer wipe(key)
	salt := make([]byte, 16)
	if _, e := rand.Read(salt); e != nil {
		return nil, nil, nil, e
	}
	k := vaultKey(password, salt)
	defer wipe(k)
	b, e := aes.NewCipher(k)
	if e != nil {
		return nil, nil, nil, e
	}
	g, e := cipher.NewGCM(b)
	if e != nil {
		return nil, nil, nil, e
	}
	nonce := make([]byte, g.NonceSize())
	if _, e = rand.Read(nonce); e != nil {
		return nil, nil, nil, e
	}
	return g.Seal(nil, nonce, key, []byte(user)), salt, nonce, nil
}
func DecryptAIKey(user string, s models.AISettings, password []byte) ([]byte, error) {
	defer wipe(password)
	k := vaultKey(password, s.Salt)
	defer wipe(k)
	b, e := aes.NewCipher(k)
	if e != nil {
		return nil, ErrAIVaultPassword
	}
	g, e := cipher.NewGCM(b)
	if e != nil || len(s.Nonce) != g.NonceSize() {
		return nil, ErrAIVaultPassword
	}
	p, e := g.Open(nil, s.Nonce, s.Ciphertext, []byte(user))
	if e != nil {
		return nil, ErrAIVaultPassword
	}
	return p, nil
}
