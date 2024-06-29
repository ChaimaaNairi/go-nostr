// ipfs_integration.go

package main

import (
	"fmt"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

var (
	ipfsAddress = "/ip4/127.0.0.1/tcp/5001"
)

// InitializeIPFSClient initializes an IPFS shell client and checks connection
func InitializeIPFSClient() *shell.Shell {
	sh := shell.NewShell(ipfsAddress)

	// Check if IPFS daemon is running
	_, err := sh.ID()
	if err != nil {
		fmt.Printf("Error connecting to IPFS daemon: %v\n", err)
		panic(err)
	}

	fmt.Println("Connected to IPFS daemon")
	return sh
}

// AddToIPFS adds a message to IPFS and returns the CID (Content ID)
func AddToIPFS(sh *shell.Shell, message string) (string, error) {
	cid, err := sh.Add(strings.NewReader(message))
	if err != nil {
		return "", fmt.Errorf("error adding message to IPFS: %v", err)
	}
	return cid, nil
}
