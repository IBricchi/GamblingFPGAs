package server

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type credential struct {
	username string
	password string
}

func AddCredential(ctx context.Context, db DB) error {
	quitStr := "#"

	fmt.Println("Adding new credentials (Enter # to quit):")
	fmt.Println("Enter username: ")
	username, quit := validateCredStrInput(quitStr)
	if quit {
		return nil
	}

	// retrieve existing data to ensure unique usernames
	creds, err := db.getCreds(ctx)
	if err != nil {
		return fmt.Errorf("server: auth: failed to get creds: %w", err)
	}
	if _, exists := creds[username]; exists {
		fmt.Println("Username already exists. Please try again.")
		return AddCredential(ctx, db)
	}

	fmt.Println("Enter password: ")
	password, quit := validateCredStrInput(quitStr)
	if quit {
		return nil
	}

	// hash with salt
	bytePassword := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("server: auth: failed to encrypt password: %w", err)
	}

	if err := db.insertCreds(ctx, credential{
		username: username,
		password: string(passwordHash),
	}); err != nil {
		return fmt.Errorf("server: auth: failed to insert credentials into db: %w", err)
	}

	// keep asking to add more credentials until quit
	return AddCredential(ctx, db)
}

// Returns a 'true' boolean and empty string if quit.
func validateCredStrInput(quitStr string) (string, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	if input == quitStr {
		return "", true
	}

	// Don't accept empty string
	if input == "" {
		fmt.Print("A value is required. Please try again: ")
		return validateCredStrInput(quitStr)
	}

	return input, false
}
