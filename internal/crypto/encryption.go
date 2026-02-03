package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// Encryptor provides AES-256-GCM encryption/decryption
type Encryptor struct {
	key []byte
}

// NewEncryptor creates a new encryptor with the given key
// Key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256
func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes, got %d", len(key))
	}
	return &Encryptor{key: key}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
// Returns ciphertext with nonce prepended
func (e *Encryptor) Encrypt(plaintext []byte) ([]byte, error) {
	if len(e.key) == 0 {
		return nil, fmt.Errorf("encryption key not set")
	}

	// Create AES cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
// Expects ciphertext with nonce prepended
func (e *Encryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(e.key) == 0 {
		return nil, fmt.Errorf("encryption key not set")
	}

	// Create AES cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: expected at least %d bytes, got %d", nonceSize, len(ciphertext))
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns ciphertext
func (e *Encryptor) EncryptString(plaintext string) ([]byte, error) {
	return e.Encrypt([]byte(plaintext))
}

// DecryptString decrypts ciphertext and returns plaintext string
func (e *Encryptor) DecryptString(ciphertext []byte) (string, error) {
	plaintext, err := e.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
