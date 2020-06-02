package plugins

//Define interface with the same set of methods as vermanCipher structure
type VermanEncryptionEngine interface {
	Encrypt(string) string
	Decrypt() (*string, error)
}

type CaesarEncryptionEngine interface {
	EncryptCaesar(shift int, text string) string
}
