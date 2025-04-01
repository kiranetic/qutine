package auth

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

const (
	salt = "qutine-salt"
)

func Authenticate() bool {
	configDir := filepath.Join(os.Getenv("HOME"), ".qutine")
	configFile := filepath.Join(configDir, "config")

	storedHash, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("qutine not installed or config missing. Run install.sh first.")
		return false
	}

	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading password")
		return false
	}
	fmt.Println()

	enteredHash := argon2.IDKey(password, []byte(salt), 1, 64*1024, 4, 32)
	
	// Debug output (remove later)
	fmt.Printf("Stored hash: %x\n", storedHash)
	fmt.Printf("Entered hash: %x\n", enteredHash)

	return bytes.Equal(enteredHash, storedHash)
}

func HashPassword(password string) []byte {
	return argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
}