package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// StrToTime convert string to time.Time with the provided layout
func StrToTime(input string, layout string) (time.Time, error) {
	r, err := time.Parse(layout, input)
	if err != nil {
		return time.Now(), err
	}
	return r, nil
}

// AskAndCompare will print prompt and compare user's input
// with expected response provided.
// Don't use this for password input.
func AskAndCompare(prompt string, expected string) (bool, error) {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	userInput = strings.TrimSuffix(userInput, "\n")
	if err != nil {
		return false, err
	}
	if userInput == expected {
		return true, nil
	}
	return false, nil
}
