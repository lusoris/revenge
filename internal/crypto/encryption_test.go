package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		wantErr bool
	}{
		{"AES-128", 16, false},
		{"AES-192", 24, false},
		{"AES-256", 32, false},
		{"Invalid size 15", 15, true},
		{"Invalid size 31", 31, true},
		{"Invalid size 0", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keySize)
			_, err := rand.Read(key)
			require.NoError(t, err)

			encryptor, err := NewEncryptor(key)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, encryptor)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, encryptor)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Generate 32-byte key for AES-256
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := NewEncryptor(key)
	require.NoError(t, err)

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{"Empty", []byte{}},
		{"Short text", []byte("hello")},
		{"Long text", []byte("The quick brown fox jumps over the lazy dog. " +
			"This is a longer message to test encryption with more data.")},
		{"Binary data", []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}},
		{"Unicode", []byte("Hello 世界 🌍")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			ciphertext, err := encryptor.Encrypt(tt.plaintext)
			require.NoError(t, err)
			assert.NotNil(t, ciphertext)

			// Ciphertext should be different from plaintext
			if len(tt.plaintext) > 0 {
				assert.NotEqual(t, tt.plaintext, ciphertext)
			}

			// Ciphertext should include nonce (12 bytes) + encrypted data + tag (16 bytes)
			assert.GreaterOrEqual(t, len(ciphertext), 12+16)

			// Decrypt
			decrypted, err := encryptor.Decrypt(ciphertext)
			require.NoError(t, err)

			// For empty input, both empty slice and nil are acceptable
			if len(tt.plaintext) == 0 {
				assert.Empty(t, decrypted)
			} else {
				assert.Equal(t, tt.plaintext, decrypted)
			}
		})
	}
}

func TestEncryptDecryptString(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := NewEncryptor(key)
	require.NoError(t, err)

	plaintext := "This is a secret message"

	// Encrypt
	ciphertext, err := encryptor.EncryptString(plaintext)
	require.NoError(t, err)

	// Decrypt
	decrypted, err := encryptor.DecryptString(ciphertext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDecrypt_DifferentNonces(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := NewEncryptor(key)
	require.NoError(t, err)

	plaintext := []byte("same plaintext")

	// Encrypt twice
	ciphertext1, err := encryptor.Encrypt(plaintext)
	require.NoError(t, err)

	ciphertext2, err := encryptor.Encrypt(plaintext)
	require.NoError(t, err)

	// Ciphertexts should be different (different nonces)
	assert.NotEqual(t, ciphertext1, ciphertext2)

	// Both should decrypt to same plaintext
	decrypted1, err := encryptor.Decrypt(ciphertext1)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted1)

	decrypted2, err := encryptor.Decrypt(ciphertext2)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted2)
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := NewEncryptor(key)
	require.NoError(t, err)

	tests := []struct {
		name       string
		ciphertext []byte
	}{
		{"Too short", []byte{0x01, 0x02}},
		{"Empty", []byte{}},
		{"Invalid data", make([]byte, 50)}, // Random data that's not encrypted
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptor.Decrypt(tt.ciphertext)
			assert.Error(t, err)
		})
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	_, err := rand.Read(key1)
	require.NoError(t, err)

	key2 := make([]byte, 32)
	_, err = rand.Read(key2)
	require.NoError(t, err)

	encryptor1, err := NewEncryptor(key1)
	require.NoError(t, err)

	encryptor2, err := NewEncryptor(key2)
	require.NoError(t, err)

	plaintext := []byte("secret message")

	// Encrypt with key1
	ciphertext, err := encryptor1.Encrypt(plaintext)
	require.NoError(t, err)

	// Try to decrypt with key2 (should fail)
	_, err = encryptor2.Decrypt(ciphertext)
	assert.Error(t, err)
}

func TestEncrypt_NoKey(t *testing.T) {
	encryptor := &Encryptor{key: []byte{}}

	_, err := encryptor.Encrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encryption key not set")
}

func TestDecrypt_NoKey(t *testing.T) {
	encryptor := &Encryptor{key: []byte{}}

	_, err := encryptor.Decrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encryption key not set")
}
