package model

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

func (c *Contract) Add(ctx TransactionContextInterface, )
