package store

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

func (c *Contract) AddPubRecord(ctx TransactionContextInterface, timestamp, deviceId string, temperature float64) error {
	return ctx.GetHotStore().AddPubRecord(&Record{
		Timestamp:   timestamp,
		DeviceID:    deviceId,
		Temperature: temperature,
	})

}

func (c *Contract) AddPvtRecord(ctx TransactionContextInterface) error {
	return ctx.GetHotStore().AddPvtRecord()
}

func (c *Contract) GetPubRecord(ctx TransactionContextInterface, timestamp, deviceId string) (*Record, error){
	record := &Record{
		Timestamp:   timestamp,
		DeviceID:    deviceId,
		Temperature: 0,
	}
	err := ctx.GetHotStore().GetPubRecord(record)
	return record, err
}

func (c *Contract) GetPvtRecord(ctx TransactionContextInterface, timestamp, deviceId string) (*Record, error) {
	record := &Record{
		Timestamp:   timestamp,
		DeviceID:    deviceId,
		Temperature: 0,
	}
	err := ctx.GetHotStore().GetPvtRecord(record)
	return record, err
}

func (c *Contract) AddArchive(ctx TransactionContextInterface) error {
	archive := &Archive{
		StartTime:     timestamp.Timestamp{},
		EndTime:       timestamp.Timestamp{},
		BlockBatchNum: "",
		Hash:          "",
	}
	return ctx.GetColdStore().AddArchive(archive)
}
