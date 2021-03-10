package store

import (
	"encoding/json"
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

func (c *Contract) AddRecord(ctx TransactionContextInterface) error {
	err := ctx.GetHotStore().AddRecord()
	if err != nil {
		return err
	}
	// TODO
	// 读取缓存目录，检查区块数量， 时机为blockNums%maxRetentionBlock+1==0, 获取当前私有数据哈希
	// 在条件达到时，对所有私有条目进行获取，记录开始-结束的主键
	// 这里需要确保主键的字典序
	return nil
}

func (c *Contract) GetRecord(ctx TransactionContextInterface, timestamp, deviceId string) ([]byte, error) {
	record := &Record{
		Timestamp: timestamp,
		DeviceID:  deviceId,
	}
	err := ctx.GetHotStore().GetRecord(record)
	if err != nil {
		return nil, err
	}
	return json.Marshal(record)
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
