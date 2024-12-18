package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

// Function to generate a random Ethereum address and private key
func generateAddressAndPrivateKey(wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()

	// Generate a random 32-byte private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Println("Error generating private key:", err)
		return
	}

	// Get the public key from the private key
	publicKey := privateKey.PublicKey

	// Derive the Ethereum address from the public key
	address := crypto.PubkeyToAddress(publicKey)

	// Convert private key to bytes (hex representation)
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// Format the address and private key as "ADDRESS | PRIVATE_KEY"
	result := fmt.Sprintf("%s | %s", address.Hex(), fmt.Sprintf("%x", privateKeyBytes))

	// Send the result to the channel
	ch <- result
}

func main() {
	// Set up concurrency
	rand.Seed(time.Now().UnixNano()) // Initialize random number generator
	var wg sync.WaitGroup
	addressCount := 1000 // Change this for how many addresses you want to generate

	// Create a channel to collect results (address | private_key)
	resultChan := make(chan string, addressCount)

	// Start address generation in parallel
	for i := 0; i < addressCount; i++ {
		wg.Add(1)
		go generateAddressAndPrivateKey(&wg, resultChan)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Write generated address | private key pairs to file
	file, err := os.Create("output/result.txt")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	// Print the results and save them
	for result := range resultChan {
		_, err := file.WriteString(result + "\n")
		if err != nil {
			log.Fatal("Error writing to file:", err)
		}
		fmt.Println(result) // Print to console
	}

	fmt.Println("Address and Private Key generation complete!")
}
