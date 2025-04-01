package main

import (
	"fmt"
	"os"

	"github.com/kiranetic/qutine/internal/auth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qtn",
	Short: "qutine - a secure, minimal container runtime",
}

var runCmd = &cobra.Command{
	Use:   "run [image] [command]",
	Short: "Run a container from an encrypted image",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !auth.Authenticate() {
			fmt.Println("Authentication failed")
			os.Exit(1)
		}
		fmt.Printf("Running container from %s with command %s\n", args[0], args[1])
	},
}

var hashCmd = &cobra.Command{
	Use:    "hash-password [password]",
	Short:  "Hash a password (internal use)",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		hash := auth.HashPassword(args[0])
		os.Stdout.Write(hash) // Output raw bytes instead of hex
	},
}

func main() {
	rootCmd.AddCommand(runCmd, hashCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}