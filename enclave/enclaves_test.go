package enclave

import (
	"os"
	"testing"
)

func createTestEnclave(key string) *Enclave {
	tempDir, err := os.MkdirTemp("", "satellion-test-*")
	if err != nil {
		panic(err)
	}
	
	enclave := &Enclave{
		decryptionKey: key,
		storagePath:   tempDir,
	}
	os.MkdirAll(tempDir, 0700)
	return enclave
}

func cleanupTestEnclave(enclave *Enclave) {
	os.RemoveAll(enclave.storagePath)
}

func TestNewEnclave(t *testing.T) {
	enclave := createTestEnclave("test-key-123")
	defer cleanupTestEnclave(enclave)
	
	// Check if storage directory was created
	if _, err := os.Stat(enclave.storagePath); os.IsNotExist(err) {
		t.Errorf("Storage directory was not created: %s", enclave.storagePath)
	}
	
	// Check if decryption key was set
	if enclave.decryptionKey != "test-key-123" {
		t.Errorf("Expected decryption key 'test-key-123', got '%s'", enclave.decryptionKey)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	enclave := createTestEnclave("test-encryption-key")
	defer cleanupTestEnclave(enclave)
	
	data := []byte("This is a test message that should be encrypted and decrypted")
	
	encryptedData, err := enclave.encrypt(data)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	
	if string(encryptedData) == string(data) {
		t.Error("Encrypted data should be different from original data")
	}
	
	decryptedData, err := enclave.decrypt(encryptedData)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}
	
	if string(decryptedData) != string(data) {
		t.Errorf("Decrypted data doesn't match original. Expected: %s, Got: %s", 
			string(data), string(decryptedData))
	}
}

func TestSaveLoad(t *testing.T) {
	enclave := createTestEnclave("test-save-load-key")
	defer cleanupTestEnclave(enclave)
	
	testData := []byte("Test data for save/load operations")
	testKey := "test-key"
	
	err := enclave.Save(testKey, testData)
	if err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}
	
	if !enclave.Exists(testKey) {
		t.Error("Saved data should exist")
	}
	
	loadedData, err := enclave.Load(testKey)
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}
	
	if string(loadedData) != string(testData) {
		t.Errorf("Loaded data doesn't match original. Expected: %s, Got: %s", 
			string(testData), string(loadedData))
	}
}

func TestLoadNonExistent(t *testing.T) {
	enclave := createTestEnclave("test-non-existent-key")
	defer cleanupTestEnclave(enclave)
	
	// Try to load non-existent data
	_, err := enclave.Load("non-existent-key")
	if err == nil {
		t.Error("Expected error when loading non-existent key")
	}
	
	// Check if it's a NotFoundError
	if !isNotFoundError(err) {
		t.Errorf("Expected NotFoundError, got: %T", err)
	}
}

func TestDelete(t *testing.T) {
	enclave := createTestEnclave("test-delete-key")
	defer cleanupTestEnclave(enclave)
	
	// Save some data first
	testData := []byte("Data to be deleted")
	testKey := "delete-test-key"
	
	err := enclave.Save(testKey, testData)
	if err != nil {
		t.Fatalf("Failed to save data for deletion test: %v", err)
	}
	
	// Verify data exists
	if !enclave.Exists(testKey) {
		t.Error("Data should exist before deletion")
	}
	
	// Delete data
	err = enclave.Delete(testKey)
	if err != nil {
		t.Fatalf("Failed to delete data: %v", err)
	}
	
	// Verify data no longer exists
	if enclave.Exists(testKey) {
		t.Error("Data should not exist after deletion")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	enclave := createTestEnclave("test-delete-non-existent")
	defer cleanupTestEnclave(enclave)
	
	// Try to delete non-existent data
	err := enclave.Delete("non-existent-key")
	if err == nil {
		t.Error("Expected error when deleting non-existent key")
	}
}

func TestList(t *testing.T) {
	enclave := createTestEnclave("test-list-key")
	defer cleanupTestEnclave(enclave)
	
	// Test initial empty list
	keys, err := enclave.List()
	if err != nil {
		t.Fatalf("Failed to list keys: %v", err)
	}
	
	if len(keys) != 0 {
		t.Errorf("Expected empty key list, got %d keys", len(keys))
	}
	
	// Save some test data
	testData1 := []byte("Test data 1")
	testData2 := []byte("Test data 2")
	
	err = enclave.Save("key1", testData1)
	if err != nil {
		t.Fatalf("Failed to save key1: %v", err)
	}
	
	err = enclave.Save("key2", testData2)
	if err != nil {
		t.Fatalf("Failed to save key2: %v", err)
	}
	
	// List keys again
	keys, err = enclave.List()
	if err != nil {
		t.Fatalf("Failed to list keys after saving: %v", err)
	}
	
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestExists(t *testing.T) {
	enclave := createTestEnclave("test-exists-key")
	defer cleanupTestEnclave(enclave)
	
	// Test non-existent key
	if enclave.Exists("non-existent") {
		t.Error("Non-existent key should not exist")
	}
	
	// Save a key
	testData := []byte("Test data")
	err := enclave.Save("test-key", testData)
	if err != nil {
		t.Fatalf("Failed to save test key: %v", err)
	}
	
	// Test existing key
	if !enclave.Exists("test-key") {
		t.Error("Existing key should exist")
	}
}

func TestEncryptDecryptWithDifferentKeys(t *testing.T) {
	enclave1 := createTestEnclave("key-1")
	defer cleanupTestEnclave(enclave1)
	enclave2 := createTestEnclave("key-2")
	defer cleanupTestEnclave(enclave2)
	
	// Test data
	originalData := []byte("Test data for different keys")
	
	// Encrypt with key 1
	encryptedData, err := enclave1.encrypt(originalData)
	if err != nil {
		t.Fatalf("Failed to encrypt with key 1: %v", err)
	}
	
	// Try to decrypt with key 2 (should fail)
	_, err = enclave2.decrypt(encryptedData)
	if err == nil {
		t.Error("Expected error when decrypting with wrong key")
	}
	
	// Decrypt with correct key 1
	decryptedData, err := enclave1.decrypt(encryptedData)
	if err != nil {
		t.Fatalf("Failed to decrypt with correct key: %v", err)
	}
	
	if string(decryptedData) != string(originalData) {
		t.Error("Decrypted data should match original data")
	}
}

func TestMultipleDataTypes(t *testing.T) {
	enclave := createTestEnclave("test-data-types")
	defer cleanupTestEnclave(enclave)
	
	// Test different data types
	testCases := []struct {
		key  string
		data []byte
	}{
		{"string-data", []byte("Hello, World!")},
		{"json-data", []byte(`{"name": "test", "value": 123}`)},
		{"binary-data", []byte{0x00, 0x01, 0x02, 0x03, 0xFF}},
		{"empty-data", []byte{}},
		{"large-data", make([]byte, 1000)}, // 1KB of zeros
	}
	
	for _, tc := range testCases {
		// Save data
		err := enclave.Save(tc.key, tc.data)
		if err != nil {
			t.Fatalf("Failed to save %s: %v", tc.key, err)
		}
		
		// Load data
		loadedData, err := enclave.Load(tc.key)
		if err != nil {
			t.Fatalf("Failed to load %s: %v", tc.key, err)
		}
		
		// Verify data matches
		if len(loadedData) != len(tc.data) {
			t.Errorf("Data length mismatch for %s: expected %d, got %d", 
				tc.key, len(tc.data), len(loadedData))
		}
		
		for i, b := range loadedData {
			if i < len(tc.data) && b != tc.data[i] {
				t.Errorf("Data mismatch for %s at position %d: expected %d, got %d", 
					tc.key, i, tc.data[i], b)
				break
			}
		}
	}
}

func TestOverwriteData(t *testing.T) {
	enclave := createTestEnclave("test-overwrite")
	defer cleanupTestEnclave(enclave)
	
	// Save initial data
	initialData := []byte("Initial data")
	err := enclave.Save("test-key", initialData)
	if err != nil {
		t.Fatalf("Failed to save initial data: %v", err)
	}
	
	// Overwrite with new data
	newData := []byte("New data")
	err = enclave.Save("test-key", newData)
	if err != nil {
		t.Fatalf("Failed to overwrite data: %v", err)
	}
	
	// Load and verify new data
	loadedData, err := enclave.Load("test-key")
	if err != nil {
		t.Fatalf("Failed to load overwritten data: %v", err)
	}
	
	if string(loadedData) != string(newData) {
		t.Errorf("Overwritten data doesn't match. Expected: %s, Got: %s", 
			string(newData), string(loadedData))
	}
}

// Helper function to check if error is NotFoundError
func isNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

