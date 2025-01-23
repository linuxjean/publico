from mnemonic import Mnemonic
from bitcoinlib.keys import HDKey
import hashlib
import base58

def private_key_to_wif(private_key_hex, compressed=True):
    try:
        # Step 1: Add version prefix (0x80 for Bitcoin mainnet)
        versioned_key = bytes.fromhex("80") + bytes.fromhex(private_key_hex)
        
        # Step 2: Add compression flag if required
        if compressed:
            versioned_key += bytes.fromhex("01")
        
        # Step 3: Double SHA-256 hash
        hash1 = hashlib.sha256(versioned_key).digest()
        hash2 = hashlib.sha256(hash1).digest()
        
        # Step 4: Append first 4 bytes of hash as checksum
        checksum = hash2[:4]
        final_key = versioned_key + checksum
        
        # Step 5: Encode in Base58
        wif_key = base58.b58encode(final_key).decode("utf-8")
        return wif_key
    except Exception as e:
        return f"Error: {e}"

# Function to generate 10 native SegWit (Bech32) addresses
def generate_native_segwit_addresses(mnemonic_words, passphrase):
    try:
        # Validate mnemonic words
        mnemo = Mnemonic("english")
        if not mnemo.check(mnemonic_words):
            raise ValueError("Invalid BIP-39 mnemonic words provided.")

        # Generate seed from mnemonic and passphrase
        seed = mnemo.to_seed(mnemonic_words, passphrase=passphrase)

        # Create an HD key object from the seed
        hd_key = HDKey.from_seed(seed, "mainnet")

        # List to store address tuples
        address_list = []

        for i in range(10):
            # Derive a child key for the path m/84'/0'/0'/0/i
            child_key = hd_key.subkey_for_path(f"m/84'/0'/0'/0/{i}")

            # Get the Bech32 address (native SegWit)
            address = child_key.address()

            # Get the public and private keys
            public_key = child_key.public_hex  # Corrected: Access as property
            private_key = child_key.private_hex  # Corrected: Access as property

            # Append the tuple to the list
            address_list.append((address, public_key, private_key))

        return address_list

    except Exception as e:
        print(f"Error: {e}")
        return []


# Main function
def main():
    print("BIP-39 Native SegWit Address Generator")

    # Input mnemonic words
    mnemonic_words = input("Enter your BIP-39 mnemonic words: ").strip()

    # Input passphrase (optional)
    passphrase = input("Enter your passphrase (optional): ").strip()

    # Generate addresses
    address_list = generate_native_segwit_addresses(mnemonic_words, passphrase)

    # Print the generated addresses
    if address_list:
        print("\nGenerated Addresses:")
        for idx, (address, pubkey, privkey) in enumerate(address_list):
            print(f"Address {idx+1}:")
            print(f"  Address: {address}")
            print(f"  Public Key: {pubkey}")
            print(f"  Private Key: {private_key_to_wif(privkey)}\n")
    else:
        print("Failed to generate addresses.")


if __name__ == "__main__":
    main()

