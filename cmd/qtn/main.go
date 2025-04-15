package main

import (
	"fmt"
	"os"
	// "path/filepath"

	"github.com/spf13/cobra"
	"github.com/kiranetic/qutine/internal/auth"
	"github.com/kiranetic/qutine/internal/container"
	"github.com/kiranetic/qutine/internal/crypto"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "qtn",
	Short: "qutine - a secure, minimal container runtime",
}

var runCmd = &cobra.Command{
	Use:   "run [image] [command]",
	Short: "Run a container from an encrypted or plain image",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !auth.Authenticate() {
			fmt.Println("Authentication failed")
			os.Exit(1)
		}

		fmt.Print("Re-enter password: ")
		pass, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		err := container.RunContainer(args[0], args[1], string(pass))
		if err != nil {
			fmt.Printf("Run failed: %v\n", err)
			os.Exit(1)
		}

		// fmt.Printf("Running container from %s with command %s\n", args[0], args[1])
	},
}

var hashCmd = &cobra.Command{
	Use:    "hash-password [password]",
	Short:  "Hash a password (internal use)",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		salt, hash := auth.GenerateSaltAndHash(args[0])
		os.Stdout.Write(append(salt, hash...)) // Output salt (16) + hash (32)
	},
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt <input.tar> <output.qimg>",
	Short: "Encrypt a container image tarball",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !auth.Authenticate() {
			fmt.Println("Authentication failed")
			os.Exit(1)
		}

		fmt.Print("Re-enter password: ")
		pass, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		err := crypto.EncryptFile(args[0], args[1], string(pass))
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Image encrypted:", args[1])
	},
}

func main() {
	rootCmd.AddCommand(runCmd, hashCmd, encryptCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
