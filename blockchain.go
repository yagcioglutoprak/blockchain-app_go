package main

import (
	"fmt"
	"net"
	"net/rpc"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"strconv"
)



// Block represents a block in the blockchain
type Block struct {
	Data      string
	Timestamp int64
	Hash      string
	PrevHash  string
	Nonce     int
}

// Blockchain represents the full blockchain
type Blockchain struct {
	Blocks []*Block
	Target string
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{Data: data, Timestamp: time.Now().Unix(), PrevHash: prevBlock.Hash}

	for newBlock.Hash = newBlock.calculateHash(); newBlock.Hash[:len(bc.Target)] != bc.Target; newBlock.Nonce++ {
		newBlock.Hash = newBlock.calculateHash()
	}

	bc.Blocks = append(bc.Blocks, newBlock)
}


// calculateHash calculates the hash of a block
func (b *Block) calculateHash() string {
	hashData := strings.Join([]string{b.Data, strconv.FormatInt(b.Timestamp, 10), b.PrevHash, strconv.Itoa(b.Nonce)}, "")
	hashInBytes := sha256.Sum256([]byte(hashData))
	return hex.EncodeToString(hashInBytes[:])
}

// RPC methods

// GetBlockchain returns the full blockchain
func (bc *Blockchain) GetBlockchain(_ struct{}, reply *[]string) error {
	hashes := make([]string, len(bc.Blocks))
	for i, b := range bc.Blocks {
		hashes[i] = b.Hash
	}
	*reply = hashes
	return nil
}

// AddBlockRPC adds a new block to the blockchain
func (bc *Blockchain) AddBlockRPC(data string, _ *struct{}) error {
	bc.AddBlock(data)
	return nil
}
func (bc *Blockchain) GetBlockDetails(hash string, reply *Block) error {
    for _, b := range bc.Blocks {
        if b.Hash == hash {
            *reply = *b
			fmt.Println(*reply)
            return nil
        }
    }
    return fmt.Errorf("Block not found")
}



func main() {
	// Create the blockchain and add the genesis block
	bc := &Blockchain{[]*Block{&Block{Data: "Genesis Block", Timestamp: time.Now().Unix(), Hash: "0"}}, "000000"}

	// Register the blockchain with the RPC server
	rpc.Register(bc)

	// Create a TCP listener
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// Accept incoming connections and handle them in a new goroutine
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
