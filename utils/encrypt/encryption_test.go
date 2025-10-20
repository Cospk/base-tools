package encrypt

import (
	"testing"
)

func TestMd5(t *testing.T) {
	// Test without salt
	expectedWithoutSalt := "098f6bcd4621d373cade4e832627b4f6" // MD5 for "test"
	resultWithoutSalt := Md5("test")
	if resultWithoutSalt != expectedWithoutSalt {
		t.Errorf("Md5 without salt = %v, want %v", resultWithoutSalt, expectedWithoutSalt)
	}

	// Test with salt
	expectedWithSalt := Md5("test", "salt") // Generate an expected value with a known salt
	if len(expectedWithSalt) == 0 {
		t.Errorf("Md5 with salt generated an empty string")
	}
}

func TestAesEncryptionDecryption(t *testing.T) {
	key := []byte("1234567890123456") // AES-128; keys are 16, 24, or 32 bytes for AES-128, AES-192, AES-256
	originalText := "Hello, World!"

	// Encrypt
	encrypted, err := AesEncrypt([]byte(originalText), key)
	if err != nil {
		t.Fatalf("AesEncrypt error: %v", err)
	}
	if len(encrypted) == 0 {
		t.Errorf("AesEncrypt returned empty data")
	}

	// Decrypt
	decrypted, err := AesDecrypt(encrypted, key)
	if err != nil {
		t.Fatalf("AesDecrypt error: %v", err)
	}
	if string(decrypted) != originalText {
		t.Errorf("AesDecrypt = %v, want %v", string(decrypted), originalText)
	}
}
