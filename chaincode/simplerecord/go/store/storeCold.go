package store

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Archive struct {
	StartTime     timestamp.Timestamp `json:"start_time"`
	EndTime       timestamp.Timestamp `json:"end_time"`
	BlockBatchNum string              `json:"block_batch_num"`
	Hash          string              `json:"hash"`
}

type TimeSeriesRecord struct {
	StartTime timestamp.Timestamp `json:"start_time"`
	EndTime   timestamp.Timestamp `json:"end_time"`
	DeviceID  string              `json:"device_id"`
	Data      []byte              `json:"data"`
}

type ColdStoreInterface interface {
	AddArchive(archive *Archive) error
	GetArchive(archive *Archive) error
}

type coldStore struct {
	Ctx contractapi.TransactionContextInterface
}

func (cs *coldStore) AddArchive(archive *Archive) error {
	archiveAsBytes, err := cs.Ctx.GetStub().GetState(archive.BlockBatchNum)
	if err != nil {
		return fmt.Errorf("Failed to get archive: " + err.Error())
	} else if archiveAsBytes != nil {
		fmt.Printf("This archive already exists: %v \n", archive.BlockBatchNum)
		return fmt.Errorf("This archive already exists: %v ", archive.BlockBatchNum)
	}

	archiveJSON, err := json.Marshal(archive)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	err = cs.Ctx.GetStub().PutState(archive.BlockBatchNum, archiveJSON)

	return err
}

func (cs *coldStore) GetArchive(archive *Archive) error {
	archiveAsBytes, err := cs.Ctx.GetStub().GetState(archive.BlockBatchNum)
	if err != nil {
		return fmt.Errorf("Failed to get archive: " + err.Error())
	}
	if archiveAsBytes == nil {
		fmt.Printf("This archive does not  exists: %v \n", archive.BlockBatchNum)
		return fmt.Errorf("This archive does not  exists: %v ", archive.BlockBatchNum)
	}
	return nil
}

func (cs *coldStore) UpdateConfig() {

}

func newColdStore(ctx contractapi.TransactionContextInterface) *coldStore {
	store := new(coldStore)
	store.Ctx = ctx
	return store
}
