package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SaveEncryptedFile encrypts and saves the uploaded file
func SaveEncryptedFile(filename string, file io.Reader) error {
	// Create /uploads directory if it doesn't exist
	uploadPath := "../uploads"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		err := os.Mkdir(uploadPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create uploads directory: %v", err)
		}
	}

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file content: %v", err)
	}

	// Encrypt the file content
	encryptedData, err := encryptData(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %v", err)
	}

	// Save encrypted file in /uploads directory
	encryptedFilePath := filepath.Join(uploadPath, filename)
	err = os.WriteFile(encryptedFilePath, encryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save encrypted file: %v", err)
	}

	return nil
}

// encryptData encrypts the data using AES encryption
func encryptData(data []byte) ([]byte, error) {
	key := []byte("obgBEboAv1av4oIU0DDiig==")
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}
