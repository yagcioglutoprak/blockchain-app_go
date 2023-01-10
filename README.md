This repository contains the code for a basic blockchain implementation with account and transfer features. The blockchain is implemented in Go and uses the RPC package to enable communication between the client and server.

The server code defines the Block and Account structs which represent blocks and account on the blockchain respectively. It also defines the Blockchain struct, which is responsible for maintaining the blockchain and handling requests from the client.

The client code allows users to interact with the blockchain by sending various HTTP requests to the server. The following endpoints are available:

/get_blockchain: Retrieve the entire blockchain from the server.

/get_block: Retrieve the details of a specific block by sending a POST request with the block's hash in the body of the request.

/add_block: Add a new block to the blockchain by sending a GET request.

/create_account: Create a new account by sending a POST request with the account's name in the body of the request.

/make_transaction: Make a transaction between two accounts by sending a POST request with the sender, receiver, and amount in the body of the request.

/get_balance: Retrieve the balance of an account by sending a POST request with the account's name in the body of the request.

To run the code, start by running the server code on your machine using the command go run blockchain.go. Next, use the client code to interact with the server by sending various HTTP requests to the server's endpoints. You can use tools such as Postman or Insomnia to test the functionality of the endpoints.

This is a basic blockchain implementation and there are many ways that it could be expanded and improved upon. The provided code is intended as a starting point and should be used as a foundation for further development.
