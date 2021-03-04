package main

import (
	"fmt"
	"record/store"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	chaincode, err := contractapi.NewChaincode(new(store.Contract))

	if err != nil {
		fmt.Printf("Error creating private record chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting private records chaincode: %s", err.Error())
	}
}
