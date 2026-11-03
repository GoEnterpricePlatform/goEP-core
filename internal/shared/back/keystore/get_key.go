package keystore

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

func (k *FileKeyStore) GetKey() ([]byte, error) {

	keyPath := filepath.Join(
		k.appDataPath,
		keyDir,
		keyFile,
	)

	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(
		strings.TrimSpace(string(data)),
	)
}
