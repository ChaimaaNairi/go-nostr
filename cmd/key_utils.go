// key_utils.go

package main

import (
	"encoding/json"
	"os"

	"github.com/nbd-wtf/go-nostr"
)

// loadKeys loads keys from a JSON file
func loadKeys(filePath string) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	var keys struct {
		PrivateKey string `json:"privateKey"`
		PublicKey  string `json:"publicKey"`
	}
	if err := json.NewDecoder(file).Decode(&keys); err != nil {
		return "", "", err
	}

	return keys.PublicKey, keys.PrivateKey, nil
}

// saveKeys saves keys to a JSON file
func saveKeys(publicKey, privateKey, filePath string) error {
	keyData := struct {
		PrivateKey string `json:"privateKey"`
		PublicKey  string `json:"publicKey"`
	}{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	jsonData, err := json.MarshalIndent(keyData, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

// generateNewKeys generates new public and private keys
func generateNewKeys() (string, string) {
	privateKey := nostr.GeneratePrivateKey()
	publicKey, _ := nostr.GetPublicKey(privateKey)
	return publicKey, privateKey
}
