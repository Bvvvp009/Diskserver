package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/scrypt"
)

var (
	ErrInvalidKey = errors.New("invalid key")
	ErrInvalidIV  = errors.New("invalid IV")
)

// Initialize these with secure, randomly generated values
var (
	encryptionKey []byte
	salt          []byte
)

func init() {
	// Generate a random encryption key and salt on startup
	encryptionKey = make([]byte, 32) // 256-bit key
	salt = make([]byte, 32)
	_, err := rand.Read(encryptionKey)
	if err != nil {
		panic("failed to generate encryption key")
	}
	_, err = rand.Read(salt)
	if err != nil {
		panic("failed to generate salt")
	}
}

func SaveEncryptedFile(filename string, content io.Reader) error {
	plaintext, err := io.ReadAll(content)
	if err != nil {
		return err
	}

	ciphertext, err := encryptFile(plaintext)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("uploads", filename), ciphertext, 0644)
}

func ReadEncryptedFile(filePath string) (io.ReadSeeker, error) {
	ciphertext, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	plaintext, err := decryptFile(ciphertext)
	if err != nil {
		return nil, err
	}

	return NewBytesReadSeeker(plaintext), nil
}

func GetFileModTime(filePath string) time.Time {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func encryptFile(plaintext []byte) ([]byte, error) {
	key, err := deriveKey(encryptionKey, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptFile(ciphertext []byte) ([]byte, error) {
	key, err := deriveKey(encryptionKey, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, ErrInvalidIV
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func deriveKey(password, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 32768, 8, 1, 32)
}

// BytesReadSeeker implements io.ReadSeeker for a byte slice
type BytesReadSeeker struct {
	data   []byte
	offset int64
}

func NewBytesReadSeeker(data []byte) *BytesReadSeeker {
	return &BytesReadSeeker{data: data}
}

func (b *BytesReadSeeker) Read(p []byte) (n int, err error) {
	if b.offset >= int64(len(b.data)) {
		return 0, io.EOF
	}
	n = copy(p, b.data[b.offset:])
	b.offset += int64(n)
	return
}

func (b *BytesReadSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.offset = offset
	case io.SeekCurrent:
		b.offset += offset
	case io.SeekEnd:
		b.offset = int64(len(b.data)) + offset
	default:
		return 0, errors.New("invalid whence")
	}

	if b.offset < 0 {
		return 0, errors.New("negative offset")
	}

	return b.offset, nil
}
