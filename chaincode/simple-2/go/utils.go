package main

import (
	"errors"
)

// GetWorldState takes the first transaction arg as the key and sets
// what is found in the world state for that key in the transaction context
func GetWorldState(ctx CustomTransactionContextInterface) error {
	_, params := ctx.GetStub().GetFunctionAndParameters()

	if len(params) < 1 {
		return errors.New("Missing key for world state")
	}

	existing, err := ctx.GetStub().GetState(params[0])

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	ctx.SetData(existing)

	return nil
}
