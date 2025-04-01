package auth

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

func Authenticate() bool {
	configDir := filepath.Join(os.Getenv("HOME"), ".qutine")
	configFile := filepath.Join(configDir, "config")

	configData, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("qutine not installed or config missing. Run install.sh first.")
		return false
	}

	// Expect config format: salt (16 bytes) + hash (32 bytes)
	if len(configData) != 16+32 {
		fmt.Println("Invalid config format")
		return false
	}
	salt := configData[:16]
	storedHash := configData[16:]

	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading password")
		return false
	}
	fmt.Println()

	enteredHash := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	return bytes.Equal(enteredHash, storedHash)
}

func HashPassword(password string, salt []byte) []byte {
	if salt == nil {
		salt = make([]byte, 16)
		_, err := rand.Read(salt)
		if err != nil {
			panic(err) // Shouldn’t happen in practice
		}
	}
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

// GenerateSaltAndHash returns salt and hash for installation
func GenerateSaltAndHash(password string) ([]byte, []byte) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return salt, hash
}