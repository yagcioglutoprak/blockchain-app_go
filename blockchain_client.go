package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"encoding/json"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Transactions []*Transaction
	Timestamp int64
	Hash      string
	PrevHash  string
	Nonce     int
}
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

// Blockchain represents the full blockchain
type Blockchain struct {
	Blocks    []*Block
	Accounts  map[string]*Account
	Target    string
}
func main() {
	// Connect to the RPC server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create a client to communicate with the server using a custom codec
	client := rpc.NewClient(conn)

	// Create an HTTP server to handle requests
	http.HandleFunc("/get_blockchain", func(w http.ResponseWriter, r *http.Request) {
		// Call the GetBlockchain method on the server to retrieve the blockchain
		var blockchain Blockchain
		err := client.Call("Blockchain.GetBlockchain", struct{}{}, &blockchain)
		if err != nil {
			http.Error(w, "Error getting blockchain", http.StatusInternalServerError)
			return
		}
		jsonResponse, _ := json.Marshal(blockchain)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

	http.HandleFunc("/get_block", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		hash := body["hash"]
		fmt.Println(hash)
		if hash == "" {
			http.Error(w, "Missing hash field", http.StatusBadRequest)
			return
		}

		// Call the GetBlockDetails function on the server
		var block Block
		if err := client.Call("Blockchain.GetBlockDetails", hash, &block); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//response
		jsonResponse, _ := json.Marshal(block)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

		

	http.HandleFunc("/add_block", func(w http.ResponseWriter, r *http.Request) {
		sender := "John"
		receiver := "Mike"
		amount := 20.5
		data := "Transfer from: " + sender + " to: " + receiver + " of Amount: " + strconv.FormatFloat(amount, 'f', -1, 64)
		var reply struct{}
	
		// Call the AddBlockRPC method on the blockchain
		err = client.Call("Blockchain.AddBlockRPC", data, &reply)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("&reply")
		}
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/create_account", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		name := body["name"]
		
		var account Account
	
		// Call the CreateAccount method on the blockchain
		err = client.Call("Blockchain.CreateAccount", name, &account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse, _ := json.Marshal(account)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

http.HandleFunc("/make_transaction", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		
		var body map[string]string
json.NewDecoder(r.Body).Decode(&body)
sender := body["sender"]
receiver := body["receiver"]
amountStr := body["amount"]
amount, _ := strconv.Atoi(amountStr)

transaction := &Transaction{sender, receiver, amount, time.Now().Unix()}
var noReply struct{}
		
		// Call the MakeTransaction method on the blockchain
		err = client.Call("Blockchain.MakeTransaction", transaction, &noReply)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/get_balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		name := body["name"]
		
		var balance int
		
		// Call the GetBalance method on the blockchain
		err = client.Call("Blockchain.GetBalance", name, &balance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		//response
		jsonResponse, _ := json.Marshal(map[string]int{"balance": balance})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})



	http.ListenAndServe(":8081", nil)
}
