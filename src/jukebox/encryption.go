package jukebox

import (
    "crypto/aes"
    "encoding/hex"
)

func EncryptAES(key []byte, plainBytes []byte) (string, error) {
    cipher, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    cipherBytes := make([]byte, len(plainBytes))
    cipher.Encrypt(cipherBytes, plainBytes)
    return hex.EncodeToString(cipherBytes), nil
}

func DecryptAES(key []byte, ecb string) ([]byte, error) {
    cipherBytes, _ := hex.DecodeString(ecb)
    cipher, err := aes.NewCipher(key)
    if err != nil {
        return []byte(""), err
    }

    plainBytes := make([]byte, len(cipherBytes))
    cipher.Decrypt(plainBytes, cipherBytes)
    return plainBytes, nil
}

