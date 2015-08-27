package otto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/scrypt"
)

const (
	cryptPrefixV0   = "v0:"
	cryptKeySaltLen = 32
)

// cryptWrite is a helper to encrypt data and then write it to a file.
// Encryption is done by using bcrypt as a KDF followed by AES-GCM.
func cryptWrite(dst string, password string, plaintext []byte) error {
	keySalt := make([]byte, cryptKeySaltLen)
	if _, err := rand.Read(keySalt); err != nil {
		return err
	}

	key, err := scrypt.Key([]byte(password), keySalt, 16384, 8, 1, 32)
	if err != nil {
		return err
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return err
	}

	// Compute random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	// Encrypt and tag with GCM
	out := gcm.Seal(nil, nonce, plaintext, nil)
	ciphertext := make([]byte, 0, len("v0:")+len(nonce)+len(out))
	ciphertext = append(ciphertext, []byte("v0:")...)
	ciphertext = append(ciphertext, keySalt...)
	ciphertext = append(ciphertext, nonce...)
	ciphertext = append(ciphertext, out...)
	out = nil

	// Create the file for writing, making sure it is opened as 0600
	// for a little additional security.
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(ciphertext))
	return err
}

func cryptRead(path string, password string) ([]byte, error) {
	// Read the contents of the path first
	ciphertext, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Verify that the data looks valid
	if !bytes.HasPrefix(ciphertext, []byte(cryptPrefixV0)) {
		return nil, fmt.Errorf("corrupt encrypted data")
	}

	// Read our key salt
	ciphertext = ciphertext[len(cryptPrefixV0):]
	keySalt := ciphertext[:32]
	ciphertext = ciphertext[32:]

	// Derive the key
	key, err := scrypt.Key([]byte(password), keySalt, 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	// Setup the cipher
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Setup the GCM AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}

	// Get the nonce and ciphertext out
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	// Decrypt
	return gcm.Open(nil, nonce, ciphertext, nil)
}
