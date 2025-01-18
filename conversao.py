import hashlib
import base58

def private_key_to_wif(private_key_hex, compressed=False):
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

# Example private key
private_key_hex = input("Digite a chave privada: ")

# Convert to WIF (uncompressed)
wif_uncompressed = private_key_to_wif(private_key_hex, compressed=False)
print("WIF (Uncompressed):", wif_uncompressed)

# Convert to WIF (compressed)
wif_compressed = private_key_to_wif(private_key_hex, compressed=True)
print("WIF (Compressed):", wif_compressed)

