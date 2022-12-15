package jukebox

import (
	"fmt"
	"testing"
)

func TestEncryptAES(t *testing.T) {
	//plainText := "Now is the time for all good men to come to the "
	plainText := "ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789"
	var key string
	var encrypted string
	var decrypted []byte
	var decryptedText string
	var err error

	// test 128-bit keys (16-bytes)
	key = "0123456789ABCDEF"
	encrypted, err = EncryptAES([]byte(key), []byte(plainText))
	if err != nil {
		t.Log("EncryptAES failed with 128-bit key")
		t.Fail()
	}

	decrypted, err = DecryptAES([]byte(key), encrypted)
	if err != nil {
		t.Log("DecryptAES failed with 128-bit key")
		t.Fail()
	}

	decryptedText = string(decrypted[:])
	if decryptedText != plainText {
		t.Log("128-bit encryption - decryptedText != plainText")
		t.Log(fmt.Sprintf("plainText = '%s'", plainText))
		t.Log(fmt.Sprintf("encrypted = '%s'", encrypted))
		t.Log(fmt.Sprintf("decrypted = '%s'", decryptedText))
		t.Fail()
	}

	// test 192-bit keys (24-bytes)
	key = "0123456789ABCDEF01234567"
	encrypted, err = EncryptAES([]byte(key), []byte(plainText))
	if err != nil {
		t.Log("EncryptAES failed with 192-bit key")
		t.Fail()
	}

	decrypted, err = DecryptAES([]byte(key), encrypted)
	if err != nil {
		t.Log("DecryptAES failed with 192-bit key")
		t.Fail()
	}

	decryptedText = string(decrypted[:])
	if decryptedText != plainText {
		t.Log("192-bit encryption - decryptedText != plainText")
		t.Log(fmt.Sprintf("plainText = '%s'", plainText))
		t.Log(fmt.Sprintf("decrypted = '%s'", decryptedText))
		t.Fail()
	}

	// test 256-bit keys (32-bytes)
	key = "0123456789ABCDEF0123456789ABCDEF"
	encrypted, err = EncryptAES([]byte(key), []byte(plainText))
	if err != nil {
		t.Log("EncryptAES failed with 256-bit key")
		t.Fail()
	}

	decrypted, err = DecryptAES([]byte(key), encrypted)
	if err != nil {
		t.Log("DecryptAES failed with 256-bit key")
		t.Fail()
	}

	decryptedText = string(decrypted[:])
	if decryptedText != plainText {
		t.Log("256-bit encryption - decryptedText != plainText")
		t.Log(fmt.Sprintf("plainText = '%s'", plainText))
		t.Log(fmt.Sprintf("decrypted = '%s'", decryptedText))
		t.Fail()
	}
}
