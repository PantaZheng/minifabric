package main

import (
	"fmt"
	"strings"
	// "recordcc/store"
	// "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	// contract := new(store.Contract)
	// contract.TransactionContextHandler = new(store.TransactionContext)
	// chaincode, err := contractapi.NewChaincode(contract)

	// if err != nil {
	// 	fmt.Printf("Error creating private record chaincode: %s", err.Error())
	// 	return
	// }

	// if err := chaincode.Start(); err != nil {
	// 	fmt.Printf("Error starting private records chaincode: %s", err.Error())
	// }
	chars := []byte(strings.Repeat("1", 1024))
	fmt.Println("---")
	fmt.Println(len(chars))
}
