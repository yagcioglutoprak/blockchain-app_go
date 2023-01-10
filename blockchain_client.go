package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"encoding/json"
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
		var hashes []string
		err := client.Call("Blockchain.GetBlockchain", struct{}{}, &hashes)
		if err != nil {
			http.Error(w, "Error getting blockchain", http.StatusInternalServerError)
			return
		}
		for _, hash := range hashes {
			fmt.Fprintln(w, hash)
		}
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

	http.ListenAndServe(":8081", nil)
}
