// ipfs_integration.go

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

var (
	// Replace with your actual IPFS daemon address
	ipfsAddress = "/ip4/127.0.0.1/tcp/5001"
)

// Initialize IPFS client
func initializeIPFSClient() *shell.Shell {
	// Connect to IPFS daemon
	sh := shell.NewShell(ipfsAddress)

	// Check if IPFS daemon is running
	_, err := sh.ID()
	if err != nil {
		fmt.Printf("Error connecting to IPFS daemon: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to IPFS daemon")

	return sh
}

// AddMessageToIPFS adds a message to IPFS and returns the IPFS hash
func AddMessageToIPFS(sh *shell.Shell, message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Add message to IPFS
	cid, err := sh.Add(ctx, shell.NewFileNode(shell.FileData([]byte(message))))
	if err != nil {
		return "", fmt.Errorf("error adding message to IPFS: %v", err)
	}

	return cid, nil
}
