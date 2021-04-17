package store

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// var data = []byte(strings.Repeat("1", 10240))

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

func (c *Contract) AddPubRecord(ctx TransactionContextInterface, timestamp int, deviceId, tmp string) error {
	if len(tmp) == 0 {
		return errors.New("tmp error")
	}
	return ctx.GetHotStore().AddPubRecord(&Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
		Data:      []byte(tmp),
	})

}

func (c *Contract) AddPvtRecord(ctx TransactionContextInterface, timestamp int, deviceId, tmp string) error {
	if len(tmp) == 0 {
		return errors.New("tmp error")
	}
	return ctx.GetHotStore().AddPvtRecord(&Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
		Data:      []byte(tmp),
	})
}

func (c *Contract) AddPubRecordTst(ctx TransactionContextInterface, timestamp int, deviceId string) error {
	return ctx.GetHotStore().AddPubRecordTst(&Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
	})
}

func (c *Contract) AddPvtRecordTst(ctx TransactionContextInterface, timestamp int, deviceId string) error {
	return ctx.GetHotStore().AddPvtRecordTst(&Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
	})
}

func (c *Contract) AddPubPvtTst(ctx TransactionContextInterface, timestamp int, deviceId string) error {
	return ctx.GetHotStore().AddPubPvtTst(&Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
	})
}

func (c *Contract) GetPubRecord(ctx TransactionContextInterface, timestamp int, deviceId string) (*Record, error) {
	record := &Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
		Data:      []byte("0"),
	}
	err := ctx.GetHotStore().GetPubRecord(record)
	return record, err
}

func (c *Contract) GetPvtRecord(ctx TransactionContextInterface, timestamp int, deviceId string) (*Record, error) {
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

func (c *Contract) HSSQ(ctx TransactionContextInterface, number int) error {
	return ctx.GetHotStore().HSSQ(number)
}

func (c *Contract) PSRQ(ctx TransactionContextInterface, number int) error {
	return ctx.GetHotStore().PSRQ(number)
}

func (c *Contract) HSMQ(ctx TransactionContextInterface, number int) error {
	return ctx.GetHotStore().HSMQ(number)
}

func (c *Contract) GS(ctx TransactionContextInterface) error {
	return ctx.GetHotStore().GS()
}

func (c *Contract) GPD(ctx TransactionContextInterface) error {
	return ctx.GetHotStore().GPD()
}

func (c *Contract) GPDH(ctx TransactionContextInterface) error {
	return ctx.GetHotStore().GPDH()
}
