/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Simple example simple Chaincode implementation
type Simple struct {
	contractapi.Contract
}

func (t *Simple) Init(ctx contractapi.TransactionContextInterface, A string, AVal int, B string, BVal int) (err error) {
	fmt.Println("Simple Init")

	fmt.Printf("AVal = %d, BVal = %d\n", AVal, BVal)

	// Write the state to the ledger
	err = ctx.GetStub().PutState(A, []byte(strconv.Itoa(AVal)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(B, []byte(strconv.Itoa(BVal)))
	if err != nil {
		return err
	}

	return err
}

// Invoke Transaction makes payment of X units from A to B
func (t *Simple) Invoke(ctx contractapi.TransactionContextInterface, A, B string, X int) (err error) {
	fmt.Println("Simple Invoke")
	var AVal, BVal int // Asset holdings

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	AValBytes, err := ctx.GetStub().GetState(A)
	if err != nil {
		return fmt.Errorf("failed to get state")
	}
	if AValBytes == nil {
		return fmt.Errorf("entity not found")
	}
	AVal, _ = strconv.Atoi(string(AValBytes))

	BValBytes, err := ctx.GetStub().GetState(B)
	if err != nil {
		return fmt.Errorf("failed to get state")
	}
	if BValBytes == nil {
		return fmt.Errorf("entity not found")
	}
	BVal, _ = strconv.Atoi(string(BValBytes))

	AVal = AVal - X
	BVal = BVal + X
	fmt.Printf("AVal = %d, BVal = %d\n", AVal, BVal)

	// Write the state back to the ledger
	err = ctx.GetStub().PutState(A, []byte(strconv.Itoa(AVal)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(B, []byte(strconv.Itoa(BVal)))
	if err != nil {
		return err
	}

	return err
}

// Delete an entity from state
func (t *Simple) Delete(ctx contractapi.TransactionContextInterface, A string) (err error) {
	fmt.Println("Simple Delete")
	// Delete the key from the state in ledger
	err = ctx.GetStub().DelState(A)
	if err != nil {
		return fmt.Errorf("failed to delete state")
	}
	return err
}

// Query callback representing the query of a chaincode
func (t *Simple) Query(ctx contractapi.TransactionContextInterface, A string) (result string, err error) {
	fmt.Println("Simple Query")
	// Get the state from the ledger
	AValBytes, err := ctx.GetStub().GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return result, errors.New(jsonResp)
	}

	if AValBytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return result, errors.New(jsonResp)
	}

	result = string(AValBytes)
	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(AValBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return result, err
}

func main() {
	cc, err := contractapi.NewChaincode(new(Simple))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting ABstore chaincode: %s", err)
	}
}
