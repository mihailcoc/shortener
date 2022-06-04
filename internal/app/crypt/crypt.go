package crypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
)

type ContextType string

const UserIDCtxName ContextType = "ctxUserId"

type Encryptor struct {
	aesblock cipher.Block
	key      []byte
}

func CookieEncryptor(key []byte) (*Encryptor, error) {
	enc := Encryptor{
		key: key,
	}
	// получаем cipher.Block
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	enc.aesblock = aesblock
	return &enc, nil
}

func (e *Encryptor) Encode(value []byte) string {
	// зашифровываем
	encrypted := make([]byte, aes.BlockSize)
	// зашифровываем
	e.aesblock.Encrypt(encrypted, value)
	// закодируем полученный массив байт
	return hex.EncodeToString(encrypted)
}

func (e *Encryptor) Decode(value string) (string, error) {
	// декодируем value в data
	encrypted, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}
	// расшифровываем
	decrypted := make([]byte, aes.BlockSize)
	// расшифровываем
	e.aesblock.Decrypt(decrypted, encrypted)
	// получаем результат в формате uuid
	result, err := uuid.FromBytes(decrypted)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func GetUserID(ctx context.Context, r *http.Request) string {
	userIDCtx := r.Context().Value(UserIDCtxName)
	userID := "default"
	if userIDCtx != nil {
		userID = userIDCtx.(string)
		return userID

	}
	return userID
}
