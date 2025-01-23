#!/usr/bin/env python3

from mnemonic import Mnemonic
import sys

def main():
    print("BIP-39 Seed Picker (Last Word Finder)")

    # Prompt user for the first 23 words
    first_23_words = input("Enter the first 23 words of your BIP-39 mnemonic, separated by spaces:\n> ").strip()
    words_23_list = first_23_words.split()

    # Make sure exactly 23 words were provided
    if len(words_23_list) != 23:
        print("\nError: You must provide exactly 23 words.")
        sys.exit(1)

    # Initialize the Mnemonic object
    mnemo = Mnemonic("english")

    # Attempt to find all possible valid last words
    valid_last_words = []
    for w in mnemo.wordlist:
        candidate = first_23_words + " " + w
        if mnemo.check(candidate):
            valid_last_words.append(w)

    # Provide feedback
    if not valid_last_words:
        print("\nNo valid last word found based on the first 23 words provided.")
    else:
        print("\nPossible valid last word(s) found:")
        for idx, w in enumerate(valid_last_words, start=1):
            print(f"{idx}. {w}")

        print("\nComplete valid mnemonic(s):")
        for w in valid_last_words:
            full_mnemonic = first_23_words + " " + w
            print(f"- {full_mnemonic}")

if __name__ == "__main__":
    main()

