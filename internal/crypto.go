package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
)

type keyFile struct {
	recipients []age.Recipient
	identities []age.Identity
}

func loadKeyFile(path string) (*keyFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return parseKeyFile(bytes.NewReader(data))
}

func parseKeyFile(r io.Reader) (*keyFile, error) {
	identities, err := age.ParseIdentities(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse identities: %w", err)
	}

	recipients := make([]age.Recipient, len(identities))
	for i, identity := range identities {
		recipients[i], err = identityToRecipient(identity)
		if err != nil {
			return nil, fmt.Errorf("failed to convert identity to recipient: %w", err)
		}
	}

	return &keyFile{recipients, identities}, nil
}

func identityToRecipient(identity age.Identity) (age.Recipient, error) {
	switch value := identity.(type) {
	case *age.HybridIdentity:
		return value.Recipient(), nil
	case *age.X25519Identity:
		return value.Recipient(), nil
	default:
		return nil, fmt.Errorf("unsupported identity type: %T", identity)
	}
}

func fetchState(dst io.Writer, keyFile *keyFile, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	r, err := age.Decrypt(file, keyFile.identities...)
	if err != nil {
		return fmt.Errorf("failed to set up decryption: %w", err)
	}

	_, err = io.Copy(dst, r)
	if err != nil {
		return fmt.Errorf("failed to copy decrypted state to destination: %w", err)
	}

	return nil
}

func updateState(src io.Reader, keyFile *keyFile, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	w, err := age.Encrypt(file, keyFile.recipients...)
	if err != nil {
		return fmt.Errorf("failed to set up encryption: %w", err)
	}
	defer w.Close()

	_, err = io.Copy(w, src)
	if err != nil {
		return fmt.Errorf("failed to encrypt into state file: %w", err)
	}

	return nil
}
