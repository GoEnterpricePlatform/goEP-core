package encryptor

type EncryptorSrv interface {
	Encrypt(value string) (string, error)
	Decrypt(value string) (string, error)
}