package store

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetHotStore() HotStoreInterface
	GetColdStore() ColdStoreInterface
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// commercial paper contract
type TransactionContext struct {
	contractapi.TransactionContext
	hotStore  *hotStore
	coldStore *coldStore
}

func (tc *TransactionContext) GetHotStore() HotStoreInterface {
	if tc.hotStore == nil {
		tc.hotStore = newHotStore(tc)
	}
	return tc.hotStore
}

func (tc *TransactionContext) GetColdStore() ColdStoreInterface {
	if tc.coldStore == nil {
		tc.coldStore = newColdStore(tc)
	}
	return tc.coldStore
}
