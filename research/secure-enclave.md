The Secure Enclave only supports NIST P-256 (secp256r1) keys and operations via Security/CryptoKit; Bitcoin requires secp256k1 (and Schnorr via BIP-340). So in-enclave signing for Bitcoin isn’t possible today.  ￼ ￼ ￼

Recommended architecture
	1.	Protect the seed at rest with Keychain + Secure Enclave gating
	•	Store an encrypted BIP39 seed (or xprv) in the Keychain.
	•	Require Face ID/Touch ID (or passcode) to unwrap using kSecAttrAccessControl (e.g., biometryCurrentSet), which is enforced by Secure Enclave.  ￼
	2.	Sign in userspace with libsecp256k1
	•	After user auth, decrypt in memory and sign PSBTs using libsecp256k1 (supports ECDSA and BIP-340 Schnorr for Taproot). Keep the key material in memory only as long as needed and zeroize.  ￼
	3.	Prefer external signers for stronger isolation
	•	Integrate hardware wallets via HWI (Ledger, Trezor, Coldcard, etc.). Your CLI can pipe PSBTs to HWI and never touch the private key.  ￼ ￼

Why Secure Enclave can’t sign Bitcoin
	•	Secure Enclave keys are EC P-256 only and created/stored inside the enclave; you can’t import a secp256k1 key, and enclave APIs don’t expose secp256k1 signing.  ￼
	•	Bitcoin uses secp256k1, and Taproot’s Schnorr signatures are BIP-340 over secp256k1.  ￼
	•	Macs with Apple silicon (and some T-chip Intel Macs) do have a Secure Enclave, just not with the curve Bitcoin needs.  ￼

Practical implementation notes for Satellion
	•	Seed generation: use the system CSPRNG (SecRandomCopyBytes) to create entropy for BIP39.  ￼
	•	Key wrapping strategy:
	•	Derive an AES key, store that key in the Keychain with biometry-gated access control, and use it to encrypt the seed/xprv you persist.  ￼
	•	Flow: PSBT → prompt user (biometry) → decrypt in RAM → sign with libsecp256k1 → wipe buffers → return PSBT.  ￼
	•	External signer option: support HWI out of the box for users who want real hardware isolation.  ￼
