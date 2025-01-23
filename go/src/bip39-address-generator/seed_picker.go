package main  // Executables must always use package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

func main() {
	fmt.Println("BIP-39 Seed Picker (Last Word Finder)")

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the first 23 words of your BIP-39 mnemonic, separated by spaces:\n> ")
	input, _ := reader.ReadString('\n')
	
	// Clean and validate input
	words23 := strings.Fields(strings.TrimSpace(input))
	if len(words23) != 23 {
		fmt.Println("\nError: You must provide exactly 23 words.")
		os.Exit(1)
	}

	// Get BIP-39 English word list
	wordList := bip39.GetWordList()
	var validLastWords []string

	// Test all possible last words
	for _, word := range wordList {
		candidate := append(words23, word)
		if bip39.IsMnemonicValid(strings.Join(candidate, " ")) {
			validLastWords = append(validLastWords, word)
		}
	}

	// Display results
	if len(validLastWords) == 0 {
		fmt.Println("\nNo valid last word found based on the first 23 words provided.")
	} else {
		fmt.Printf("\nPossible valid last word(s) found (%d):\n", len(validLastWords))
		for i, w := range validLastWords {
			fmt.Printf("%d. %s\n", i+1, w)
		}

		fmt.Println("\nComplete valid mnemonic(s):")
		basePhrase := strings.Join(words23, " ")
		for _, w := range validLastWords {
			fmt.Printf("- %s %s\n", basePhrase, w)
		}
	}
}
