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

type Account struct {
	Name       string
	Balance    int
	Transactions []*Transaction
}

type Transaction struct {
	Sender    string
	Recipient string
	Amount    int
	Timestamp int64
}

// Block represents a block in the blockchain
type Block struct {
	Transactions []*Transaction
	Timestamp int64
	Hash      string
	PrevHash  string
	Nonce     int
}

// Blockchain represents the full blockchain
type Blockchain struct {
	Blocks    []*Block
	Accounts  map[string]*Account
	Target    string
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{Transactions: transactions, Timestamp: time.Now().Unix(), PrevHash: prevBlock.Hash}

	for newBlock.Hash = newBlock.calculateHash(); newBlock.Hash[:len(bc.Target)] != bc.Target; newBlock.Nonce++ {
		newBlock.Hash = newBlock.calculateHash()
	}

	bc.Blocks = append(bc.Blocks, newBlock)
}

// CreateAccount creates a new account
func (bc *Blockchain) CreateAccount(name string, reply *Account) error {
	if _, ok := bc.Accounts[name]; ok {
		return fmt.Errorf("Account with name %s already exists", name)
	}
	newAccount := &Account{Name: name, Balance: 100}
	bc.Accounts[name] = newAccount
	*reply = *newAccount
	fmt.Println(newAccount)
	return nil
}

// MakeTransaction makes a transaction between two accounts
func (bc *Blockchain) MakeTransaction(transaction *Transaction, _ *struct{}) error {
	sender := bc.Accounts[transaction.Sender]
	if sender == nil {
		return fmt.Errorf("Sender account not found")
	}
	recipient := bc.Accounts[transaction.Recipient]
	if recipient == nil {
		return fmt.Errorf("Recipient account not found")
	}
	if sender.Balance < transaction.Amount {
		return fmt.Errorf("Sender has insufficient funds")
	}

	sender.Balance -= transaction.Amount
	sender.Transactions = append(sender.Transactions, transaction)
	recipient.Balance += transaction.Amount
	recipient.Transactions = append(recipient.Transactions, transaction)
	bc.AddBlock([]*Transaction{transaction})
	return nil
}

// GetBalance gets the balance of an account
func (bc *Blockchain) GetBalance(name string, reply *int) error {
	account := bc.Accounts[name]
	if account == nil {
		return fmt.Errorf("Account not found")
	}
	*reply = account.Balance
	return nil
}

// calculateHash calculates the hash of a block
func (b *Block) calculateHash() string {
	transactionsData := make([]string, len(b.Transactions))
	for i, t := range b.Transactions {
		transactionsData[i] = t.Sender + t.Recipient + strconv.Itoa(t.Amount) + strconv.FormatInt(t.Timestamp, 10)
	}
	hashData := strings.Join(append(transactionsData, strconv.FormatInt(b.Timestamp, 10), b.PrevHash, strconv.Itoa(b.Nonce)), "")
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

func (bc *Blockchain) GetBlockDetails(hash string, reply *Block) error {
    for _, b := range bc.Blocks {
        if b.Hash == hash {
            *reply = *b
            return nil
        }
    }
    return fmt.Errorf("Block not found")
}


func (bc *Blockchain) GetAccountDetails(name string, reply *Account) error {
    account := bc.Accounts[name]
    if account == nil {
        return fmt.Errorf("Account not found")
    }
    *reply = *account
    return nil
}

func (bc *Blockchain) GetBlockTransactions(hash string, reply *[]*Transaction) error {
    for _, b := range bc.Blocks {
        if b.Hash == hash {
            *reply = b.Transactions
            return nil
        }
    }
    return fmt.Errorf("Block not found")
}

func main() {
	// Create the blockchain and add the genesis block
	bc := &Blockchain{
		Blocks:    []*Block{&Block{Transactions: []*Transaction{}, Timestamp: time.Now().Unix(), Hash: "0"}},
		Accounts:  map[string]*Account{},
		Target:    "000000",
	}

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