package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
)

type AddressInfo struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

func privateKeyToWIF(privateKey []byte, compressed bool) (string, error) {
	versionedKey := append([]byte{0x80}, privateKey...)
	if compressed {
		versionedKey = append(versionedKey, 0x01)
	}

	firstSHA := sha256.Sum256(versionedKey)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	finalKey := append(versionedKey, checksum...)
	return base58.Encode(finalKey), nil
}

func generateNativeSegwitAddresses(mnemonic, passphrase string) ([]AddressInfo, error) {
	var addressList []AddressInfo

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid BIP-39 mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, passphrase)
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %v", err)
	}

	currentKey := masterKey
	for _, index := range []uint32{
		hdkeychain.HardenedKeyStart + 84,
		hdkeychain.HardenedKeyStart + 0,
		hdkeychain.HardenedKeyStart + 0,
	} {
		currentKey, err = currentKey.Derive(index)
		if err != nil {
			return nil, fmt.Errorf("derivation failed: %v", err)
		}
	}

	currentKey, err = currentKey.Derive(0)
	if err != nil {
		return nil, fmt.Errorf("external chain derivation failed: %v", err)
	}

	for i := 0; i < 10; i++ {
		childKey, err := currentKey.Derive(uint32(i))
		if err != nil {
			fmt.Printf("Child %d derivation failed: %v\n", i, err)
			continue
		}

		privKey, err := childKey.ECPrivKey()
		if err != nil {
			fmt.Printf("Private key extraction failed for %d: %v\n", i, err)
			continue
		}

		privKeyBytes := privKey.Serialize()
		wif, err := privateKeyToWIF(privKeyBytes, true)
		if err != nil {
			fmt.Printf("WIF conversion failed for %d: %v\n", i, err)
			continue
		}

		pubKeyBytes := privKey.PubKey().SerializeCompressed()
		pubKeyHash := btcutil.Hash160(pubKeyBytes)
		addressWitness, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
		if err != nil {
			fmt.Printf("Address creation failed for %d: %v\n", i, err)
			continue
		}

		addressList = append(addressList, AddressInfo{
			Address:    addressWitness.EncodeAddress(),
			PublicKey:  hex.EncodeToString(pubKeyBytes),
			PrivateKey: wif,
		})
	}

	return addressList, nil
}

func main() {
	fmt.Println("BIP-39 Native SegWit Address Generator")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your BIP-39 mnemonic words: ")
	mnemonic, _ := reader.ReadString('\n')
	mnemonic = strings.TrimSpace(mnemonic)

	fmt.Print("Enter your passphrase (optional): ")
	passphrase, _ := reader.ReadString('\n')
	passphrase = strings.TrimSpace(passphrase)

	addressList, err := generateNativeSegwitAddresses(mnemonic, passphrase)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(addressList) > 0 {
		fmt.Println("\nGenerated Addresses:")
		for i, addr := range addressList {
			fmt.Printf("Address %d:\n", i+1)
			fmt.Printf("  Address: %s\n", addr.Address)
			fmt.Printf("  Public Key: %s\n", addr.PublicKey)
			fmt.Printf("  Private Key: %s\n\n", addr.PrivateKey)
		}
	} else {
		fmt.Println("No addresses generated")
	}
}
