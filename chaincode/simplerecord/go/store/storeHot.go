package store

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const collectionName = "collectionPrivate"

//var emptyValue []byte

//type RetentionPolicy struct {
//	StartTime timestamp.Timestamp `json:"start_time"`
//	Interval  int64               `json:"interval"`
//}

type Record struct {
	Timestamp string `json:"timestamp"`
	DeviceID  string `json:"device_id"`
	Data      []byte `json:"data"`
}

func (r *Record) makePrimaryKey() string {
	return strings.Join([]string{r.Timestamp, r.DeviceID}, "_")
}

type HotStoreInterface interface {
	AddPubRecord(record *Record) error
	AddPvtRecord(record *Record) error
	GetPubRecord(record *Record) error
	GetPvtRecord(record *Record) error
	//GetRecordHash(timestamp, id string) (string, error)
	//GetRecordsByRange(startKey, endKey string) ([]Record, error)
}

type hotStore struct {
	Ctx contractapi.TransactionContextInterface
}

func (hs *hotStore) AddPubRecord(record *Record) error {
	primaryKey := record.makePrimaryKey()
	//// 检查重复
	//recordAsBytes, err := hs.Ctx.GetStub().GetState(primaryKey)
	//if err != nil {
	//	return fmt.Errorf("Failed to get record: " + err.Error())
	//} else if recordAsBytes != nil {
	//	fmt.Println("This record already exists: " + primaryKey)
	//	return fmt.Errorf("This record already exists: " + primaryKey)
	//}

	hotRecordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// 保存热点记录到状态库
	err = hs.Ctx.GetStub().PutState(primaryKey, hotRecordJSON)
	if err != nil {
		return fmt.Errorf("failed to put hot record: %s", err.Error())
	}

	// 时序索引
	/*indexName := "timestamp~id"
	timeIndex, err := hs.Ctx.GetStub().CreateCompositeKey(indexName, []string{record.Timestamp, record.DeviceID})
	if err != nil {
		return err
	}
	err = hs.Ctx.GetStub().PutState(timeIndex, emptyValue)*/

	return err
}

func (hs *hotStore) AddPvtRecord(record *Record) error {
	//transMap, err := hs.Ctx.GetStub().GetTransient()
	//if err != nil {
	//	return fmt.Errorf("Error getting transient: " + err.Error())
	//}
	//
	//transientRecordJSON, ok := transMap["record"]
	//if !ok {
	//	return fmt.Errorf("hot_record not found in the transient map")
	//}
	//
	//var recordInput Record
	//err = json.Unmarshal(transientRecordJSON, &recordInput)
	//if err != nil {
	//	return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	//}
	//
	//record := &Record{
	//	Timestamp:   recordInput.Timestamp,
	//	DeviceID:    recordInput.DeviceID,
	//	Temperature: recordInput.Temperature,
	//}

	primaryKey := record.makePrimaryKey()

	//// 检查重复
	//recordAsBytes, err := hs.Ctx.GetStub().GetPrivateData(collectionName, primaryKey)
	//if err != nil {
	//	return fmt.Errorf("Failed to get record: " + err.Error())
	//} else if recordAsBytes != nil {
	//	fmt.Println("This record already exists: " + primaryKey)
	//	return fmt.Errorf("This record already exists: " + primaryKey)
	//}

	hotRecordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// 保存热点记录到状态库
	err = hs.Ctx.GetStub().PutPrivateData(collectionName, primaryKey, hotRecordJSON)
	if err != nil {
		return fmt.Errorf("failed to put hot record: %s", err.Error())
	}

	// 时序索引
	//indexName := "timestamp~id"
	//timeIndex, err := hs.Ctx.GetStub().CreateCompositeKey(indexName, []string{record.Timestamp, record.DeviceID})
	//if err != nil {
	//	return err
	//}
	//err = hs.Ctx.GetStub().PutPrivateData(collectionName, timeIndex, emptyValue)

	return err
}

func (hs *hotStore) GetPubRecord(record *Record) error {
	primaryKey := record.makePrimaryKey()
	recordAsBytes, err := hs.Ctx.GetStub().GetState(primaryKey)
	if err != nil {
		return fmt.Errorf("failed to read from marble %s", err.Error())
	}
	if recordAsBytes == nil {
		return fmt.Errorf("%s does not exist", primaryKey)
	}
	err = json.Unmarshal(recordAsBytes, record)
	return err
}

func (hs *hotStore) GetPvtRecord(record *Record) error {
	primaryKey := record.makePrimaryKey()
	recordAsBytes, err := hs.Ctx.GetStub().GetPrivateData(collectionName, primaryKey)
	if err != nil {
		return fmt.Errorf("failed to read from marble %s", err.Error())
	}
	if recordAsBytes == nil {
		return fmt.Errorf("%s does not exist", primaryKey)
	}
	err = json.Unmarshal(recordAsBytes, record)
	return err
}

//func (hs *hotStore) GetRecordHash(timestamp, id string) (string, error) {
//	r := Record{
//		Timestamp: timestamp,
//		DeviceID:        id,
//	}
//	primaryKey := r.makePrimaryKey()
//	hashAsBytes, err := hs.Ctx.GetStub().GetPrivateDataHash(collectionName, primaryKey)
//	if err != nil {
//		return "", fmt.Errorf("Failed to get public data hash for record:" + err.Error())
//	} else if hashAsBytes == nil {
//		return "", fmt.Errorf("Record does not exist: " + primaryKey)
//	}
//
//	return string(hashAsBytes), nil
//}

//func (hs *hotStore) GetRecordsByRange(startKey, endKey string) ([]Record, error) {
//	iterator, err := hs.Ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		_ = iterator.Close()
//	}()
//
//	var records []Record
//
//	for iterator.HasNext() {
//		currentItem, err := iterator.Next()
//		if err != nil {
//			return nil, err
//		}
//
//		newRecord := new(Record)
//		err = json.Unmarshal(currentItem.Value, newRecord)
//		if err != nil {
//			return nil, err
//		}
//
//		records = append(records, *newRecord)
//	}
//
//	return records, nil
//}

func newHotStore(ctx contractapi.TransactionContextInterface) *hotStore {
	store := new(hotStore)
	store.Ctx = ctx
	return store
}
