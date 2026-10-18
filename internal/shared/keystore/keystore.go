package keystore

type KeyStore interface {
	GetKey() ([]byte, error)
	CreateKey() ([]byte, error)
	LoadOrCreateKey() ([]byte, error)
}