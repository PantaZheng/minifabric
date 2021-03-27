package store

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var data = []byte(strings.Repeat("1", 10240))

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

func (c *Contract) AddPubRecord(ctx TransactionContextInterface, timestamp, deviceId string) error {
	return ctx.GetHotStore().AddPubRecord(&Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
		Data:      data,
	})

}

func (c *Contract) AddPvtRecord(ctx TransactionContextInterface, timestamp, deviceId string) error {
	return ctx.GetHotStore().AddPvtRecord(&Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
		Data:      data,
	})
}

func (c *Contract) GetPubRecord(ctx TransactionContextInterface, timestamp, deviceId string) (*Record, error) {
	record := &Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
		Data:      []byte("0"),
	}
	err := ctx.GetHotStore().GetPubRecord(record)
	return record, err
}

func (c *Contract) GetPvtRecord(ctx TransactionContextInterface, timestamp, deviceId string) (*Record, error) {
	record := &Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
		Data:      []byte("0"),
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
