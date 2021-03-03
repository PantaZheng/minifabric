package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"marbles/model"
)

func main() {

	chaincode, err := contractapi.NewChaincode(new(model.SmartContract))

	if err != nil {
		fmt.Printf("Error creating private record chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting private records chaincode: %s", err.Error())
	}
}
