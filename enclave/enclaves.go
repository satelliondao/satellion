package enclave

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Enclave struct {
	decryptionKey string
	storagePath   string
}

func NewEnclave(decryptionKey string) *Enclave {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	storagePath := filepath.Join(homeDir, ".satellion", "enclave")
	
	os.MkdirAll(storagePath, 0700)
	return &Enclave{
		decryptionKey: decryptionKey,
		storagePath:   storagePath,
	}
}

// deriveKey creates a 32-byte key from the decryption key using SHA-256
func (e *Enclave) deriveKey() []byte {
	hash := sha256.Sum256([]byte(e.decryptionKey))
	return hash[:]
}

// generateNonce creates a random 12-byte nonce for AES-GCM
func (e *Enclave) generateNonce() ([]byte, error) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return nonce, nil
}

func (e *Enclave) encrypt(data []byte) ([]byte, error) {
	key := e.deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce, err := e.generateNonce()
	if err != nil {
		return nil, err
	}

	// Encrypt and append nonce
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (e *Enclave) decrypt(encryptedData []byte) ([]byte, error) {
	key := e.deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(encryptedData) < 12 {
		return nil, fmt.Errorf("encrypted data too short")
	}

	// Extract nonce and ciphertext
	nonce := encryptedData[:12]
	ciphertext := encryptedData[12:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (e *Enclave) getFilePath(key string) string {
	hash := sha256.Sum256([]byte(key))
	filename := hex.EncodeToString(hash[:]) + ".enc"
	return filepath.Join(e.storagePath, filename)
}

func (e *Enclave) Save(key string, data []byte) error {
	encryptedData, err := e.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	filePath := e.getFilePath(key)
	err = os.WriteFile(filePath, encryptedData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return nil
}

func (e *Enclave) Load(key string) ([]byte, error) {
	filePath := e.getFilePath(key)
	
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &NotFoundError{Key: key}
		}
		return nil, fmt.Errorf("failed to read encrypted file: %w", err)
	}

	decryptedData, err := e.decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decryptedData, nil
}

func (e *Enclave) Delete(key string) error {
	filePath := e.getFilePath(key)
	
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("key '%s' not found in storage", key)
		}
		return fmt.Errorf("failed to delete encrypted file: %w", err)
	}

	return nil
}

func (e *Enclave) List() ([]string, error) {
	files, err := os.ReadDir(e.storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory: %w", err)
	}

	var keys []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".enc" {
			// For now, we can't easily reverse the hash to get the original key
			// This is a limitation of the current design
			keys = append(keys, file.Name())
		}
	}

	return keys, nil
}

func (e *Enclave) Exists(key string) bool {
	filePath := e.getFilePath(key)
	_, err := os.Stat(filePath)
	return err == nil
}

func (e *Enclave) GetStoragePath() string {
	return e.storagePath
}

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("key '%s' not found in storage", e.Key)
}