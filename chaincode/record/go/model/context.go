package model

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"marbles/store"
)

type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetHotStore() store.HotStoreInterface
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// commercial paper contract
type TransactionContext struct {
	contractapi.TransactionContext
	coldStore *store.coldStore
	hotStore  *store.hotStore
}

func (tc *TransactionContext) GetHotStore() store.HotStoreInterface {
	if tc.hotStore == nil {
		tc.hotStore = store.newHotStore(tc)
	}

	return tc.hotStore
}
