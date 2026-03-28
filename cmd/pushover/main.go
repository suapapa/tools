package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var title, msg string
	switch len(os.Args) {
	case 2:
		msg = os.Args[1]
	case 3:
		title = os.Args[1]
		msg = os.Args[2]
	default:
		printUsage()
		return nil
	}

	userKey := os.Getenv("PUSHOVER_USERKEY")
	appToken := os.Getenv("PUSHOVER_APPTOKEN")
	if userKey == "" || appToken == "" {
		return fmt.Errorf("PUSHOVER_USERKEY and PUSHOVER_APPTOKEN must be set")
	}

	formValues := url.Values{
		"token":   {appToken},
		"user":    {userKey},
		"title":   {title},
		"message": {msg},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.pushover.net/1/messages.json", strings.NewReader(formValues.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	fmt.Println()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("  %s [title] message\n", os.Args[0])
	fmt.Println("\nEnvironment variables:")
	fmt.Println("  PUSHOVER_USERKEY")
	fmt.Println("  PUSHOVER_APPTOKEN")
}
