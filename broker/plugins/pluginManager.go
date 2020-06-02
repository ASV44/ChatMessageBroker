package plugins

import (
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"plugin"
)

type PluginManager struct {

}

func (manager PluginManager) loadCipherPlugin() (EncryptionPluginEngine, error) {
	// Load Cipher plugin
	pluginModule, err := plugin.Open("./encrypting/cipher.so")
	if err != nil {
		return EncryptionPluginEngine{}, models.PluginError{
			Message: "Unable to load cipher module",
			Err: err,
		}
	}

	//Load DecryptCaesar function
	caesarCipherSymbol, err := pluginModule.Lookup("CaesarCipher")
	if err != nil {
		return EncryptionPluginEngine{}, models.PluginError{
			Message: "Unable to load caesar decrypt function",
			Err: err,
		}
	}

	//Load VermanCipher variable
	vermanCipherSymbol, err := pluginModule.Lookup("VermanCipher")
	if err != nil {
		return EncryptionPluginEngine{}, models.PluginError{
			Message: "Unable to load VermanCipher variable",
			Err: err,
		}
	}

	//Cast encryptCaesar symbol to the correct type
	encryptCaesar := caesarCipherSymbol.(CaesarEncryptionEngine)
	//Cast vermanCipher symbol to the correct interface type
	vermanCipherIf := vermanCipherSymbol.(VermanEncryptionEngine)

	return EncryptionPluginEngine{Caeser: encryptCaesar, Verman: vermanCipherIf}, nil
}
