package store

import (
	"encoding/json"
	"fmt"
	"marbles/model"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const collectionName = "hot_records"

var emptyValue = []byte{0x00}

type Record struct {
	ID          string  `json:"device_id"`
	Timestamp   string  `json:"timestamp"`
	Temperature float64 `json:"temperature"`
}

type HotStoreInterface interface {
	AddRecord() error
	GetRecord(timestamp string, id string) (*Record, error)
	GetRecordHash(timestamp string, id string) (string, error)
	GetRecordsByRange(startKey string, endKey string) ([]Record, error)
}

type hotStore struct {
	Ctx contractapi.TransactionContextInterface
}

func (hs *hotStore) makePrimaryKey(timestamp string, id string) string {
	return strings.Join([]string{timestamp, id}, "_")
}

func (hs *hotStore) AddRecord() error {
	transMap, err := hs.Ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	transientRecordJSON, ok := transMap["record"]
	if !ok {
		return fmt.Errorf("hot_record not found in the transient map")
	}

	var recordInput Record
	err = json.Unmarshal(transientRecordJSON, &recordInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(recordInput.ID) == 0 {
		return fmt.Errorf("id field must be a non-empty string")
	}

	if len(recordInput.Timestamp) == 0 {
		return fmt.Errorf("timestamp field must be a non-empty string")
	}

	if recordInput.Temperature < 0 || recordInput.Temperature > 100 {
		return fmt.Errorf("temperature field must between 0 ~ 100")
	}

	primaryKey := hs.makePrimaryKey(recordInput.Timestamp, recordInput.ID)
	// 检查重复
	recordAsBytes, err := hs.Ctx.GetStub().GetPrivateData(collectionName, primaryKey)
	if err != nil {
		return fmt.Errorf("Failed to get recordInput: " + err.Error())
	} else if recordAsBytes != nil {
		fmt.Println("This recordInput already exists: " + primaryKey)
		return fmt.Errorf("This recordInput already exists: " + primaryKey)
	}

	hotRecord := &Record{
		ID:          recordInput.ID,
		Timestamp:   recordInput.Timestamp,
		Temperature: recordInput.Temperature,
	}

	hotRecordJSON, err := json.Marshal(hotRecord)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// 保存热点记录到状态库
	err = hs.Ctx.GetStub().PutPrivateData(collectionName, primaryKey, hotRecordJSON)
	if err != nil {
		return fmt.Errorf("failed to put hot record: %s", err.Error())
	}

	// 时序索引
	indexName := "timestamp~id"
	timeIndex, err := hs.Ctx.GetStub().CreateCompositeKey(indexName, []string{recordInput.Timestamp, recordInput.ID})
	if err != nil {
		return err
	}
	err = hs.Ctx.GetStub().PutPrivateData(collectionName, timeIndex, emptyValue)

	return err
}

func (hs *hotStore) GetRecord(timestamp string, id string) (*Record, error) {
	primaryKey := hs.makePrimaryKey(timestamp, id)
	recordAsBytes, err := hs.Ctx.GetStub().GetPrivateData(collectionName, primaryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from marble %s", err.Error())
	}
	if recordAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", primaryKey)
	}

	record := new(Record)
	err = json.Unmarshal(recordAsBytes, record)

	return record, err
}

func (hs *hotStore) GetRecordHash(timestamp string, id string) (string, error) {
	primaryKey := hs.makePrimaryKey(timestamp, id)
	hashAsBytes, err := hs.Ctx.GetStub().GetPrivateDataHash(collectionName, primaryKey)
	if err != nil {
		return "", fmt.Errorf("Failed to get public data hash for record:" + err.Error())
	} else if hashAsBytes == nil {
		return "", fmt.Errorf("Record does not exist: " + primaryKey)
	}

	return string(hashAsBytes), nil
}

func (hs *hotStore) GetRecordsByRange(startKey string, endKey string) ([]Record, error) {
	iterator, err := hs.Ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = iterator.Close()
	}()

	var records []Record

	for iterator.HasNext() {
		currentItem, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		newRecord := new(Record)
		err = json.Unmarshal(currentItem.Value, newRecord)
		if err != nil {
			return nil, err
		}

		records = append(records, *newRecord)
	}

	return records, nil
}

func newHotStore(ctx model.TransactionContextInterface) *hotStore {
	store := new(hotStore)
	store.Ctx = ctx
	return store
}
