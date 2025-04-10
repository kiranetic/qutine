package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
	"io"
	"os"
)

// DeriveKeyFromPassword returns a 32-byte key from password and salt.
func DeriveKeyFromPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

// EncryptFile encrypts the input file using AES-CFB and saves to outPath.
// File format: [salt (16B)] [iv (16B)] [encrypted_data]
func EncryptFile(inPath, outPath string, password string) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := DeriveKeyFromPassword(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write salt and IV
	if _, err := outFile.Write(salt); err != nil {
		return err
	}
	if _, err := outFile.Write(iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outFile}

	_, err = io.Copy(writer, inFile)
	return err
}

// DecryptFile decrypts a .qimg file encrypted with EncryptFile.
func DecryptFile(inPath, outPath string, password string) error {
	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	header := make([]byte, 32) // 16 bytes salt + 16 bytes IV
	if _, err := io.ReadFull(inFile, header); err != nil {
		return err
	}
	salt := header[:16]
	iv := header[16:]

	key := DeriveKeyFromPassword(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inFile}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	return err
}
