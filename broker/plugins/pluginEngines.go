package plugins

//Define interface with the same set of methods as vermanCipher structure
type EncryptionEngine interface {
	Encrypt(string) string
	Decrypt() (*string, error)
}
