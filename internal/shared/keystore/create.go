package keystore

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"
)

func (k *FileKeyStore) CreateKey() ([]byte, error) {

	keyPath := filepath.Join(
		k.appDataPath,
		keyDir,
		keyFile,
	)

	dir := filepath.Dir(keyPath)

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return nil, err
	}

	key := make([]byte, keySize)

	_, err = rand.Read(key)
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(key)

	err = os.WriteFile(
		keyPath,
		[]byte(encoded),
		0600,
	)
	if err != nil {
		return nil, err
	}

	return key, nil
}
