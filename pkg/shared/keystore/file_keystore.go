package keystore

import (
	"errors"
	"os"
)

const (
	keySize = 32
	keyDir  = "keys"
	keyFile = "payment-encryption.key"
)

type FileKeyStore struct {
	appDataPath string
}

func NewFileKeyStore(appDataPath string) *FileKeyStore {
	return &FileKeyStore{
		appDataPath: appDataPath,
	}
}

func (k *FileKeyStore) LoadOrCreateKey() ([]byte, error) {

	key, err := k.GetKey()
	if err == nil {
		return key, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return k.CreateKey()
	}

	return nil, err
}