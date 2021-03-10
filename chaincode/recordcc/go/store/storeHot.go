package store

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const collectionName = "collectionPrivate"

var emptyValue []byte

//type RetentionPolicy struct {
//	StartTime timestamp.Timestamp `json:"start_time"`
//	Interval  int64               `json:"interval"`
//}

type Point struct {
	Timestamp   timestamp.Timestamp `json:"timestamp"`
	Temperature float64             `json:"temperature"`
}

type Record struct {
	Timestamp timestamp.Timestamp `json:"timestamp"`
	Points    []Point             `json:"points"`
	ClientID  string              `json:"client_id"`
}

func (r *Record) makePrimaryKey() string {
	return strings.Join([]string{r.Timestamp, r.ClientID}, "_")
}

type HotStoreInterface interface {
	AddPvtRecord() error
	GetPvtRecord(record *Record) error
	GetPvtRecordHash(record *Record) (string, error)
	//GetRecordsByRange(startKey, endKey string) ([]Record, error)
}

type recordStore struct {
	Ctx    contractapi.TransactionContextInterface
	Points []Point
}

//func (hs *recordStore) archive() error {
//	files, _ := ioutil.ReadDir("./")
//}

func (rs *recordStore) AddPoint() error {
	peer_id, err := rs.Ctx.GetClientIdentity().GetID()

}

func (rs *recordStore) newRecord() error {

}

func (rs *recordStore) AddPvtRecord() error {
	transMap, err := rs.Ctx.GetStub().GetTransient()
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

	if len(recordInput.UUID) == 0 {
		return fmt.Errorf("uuid field must be a non-empty string")
	}

	primaryKey := recordInput.makePrimaryKey()
	// 检查重复
	recordAsBytes, err := rs.Ctx.GetStub().GetPrivateData(collectionName, primaryKey)
	if err != nil {
		return fmt.Errorf("Failed to get recordInput: " + err.Error())
	} else if recordAsBytes != nil {
		fmt.Println("This recordInput already exists: " + primaryKey)
		return fmt.Errorf("This recordInput already exists: " + primaryKey)
	}

	hotRecord := &Record{
		UUID:      recordInput.UUID,
		Timestamp: recordInput.Timestamp,
		Data:      recordInput.Data,
	}

	hotRecordJSON, err := json.Marshal(hotRecord)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// 保存热点记录到状态库
	err = rs.Ctx.GetStub().PutPrivateData(collectionName, primaryKey, hotRecordJSON)
	if err != nil {
		return fmt.Errorf("failed to put hot record: %s", err.Error())
	}

	// timestamp
	indexName := "uuid"
	timeIndex, err := rs.Ctx.GetStub().CreateCompositeKey(indexName, []string{recordInput.UUID})
	if err != nil {
		return err
	}
	err = rs.Ctx.GetStub().PutPrivateData(collectionName, timeIndex, emptyValue)

	fmt.Println("recordInput:", recordInput)

	return err
}

func (rs *recordStore) GetPvtRecord(record *Record) error {
	primaryKey := record.makePrimaryKey()
	recordAsBytes, err := rs.Ctx.GetStub().GetPrivateData(collectionName, primaryKey)
	if err != nil {
		return fmt.Errorf("failed to read from marble %s", err.Error())
	}
	if recordAsBytes == nil {
		return fmt.Errorf("%s does not exist", primaryKey)
	}

	err = json.Unmarshal(recordAsBytes, record)
	return err
}

func (rs *recordStore) GetPvtRecordHash(record *Record) (string, error) {
	primaryKey := record.makePrimaryKey()
	hashAsBytes, err := rs.Ctx.GetStub().GetPrivateDataHash(collectionName, primaryKey)
	if err != nil {
		return "", fmt.Errorf("Failed to get private data hash for record:" + err.Error())
	} else if hashAsBytes == nil {
		return "", fmt.Errorf("Record does not exist: " + primaryKey)
	}

	return string(hashAsBytes), nil
}

func (rs *recordStore) GetRecordsByRange(startKey, endKey string) ([]Record, error) {
	iterator, err := rs.Ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)
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

func newHotStore(ctx contractapi.TransactionContextInterface) *recordStore {
	//获取首次运行时间
	//将上传数据条数与区块号传入公有数据，存入哈希值
	//将时间片，节点编号，哈希，串写入私有数据
	//创建一个数据缓冲区间，时序数据，自动时间归档，存入
	//在每次交易执行完之后进行判断
	store := new(recordStore)
	store.Ctx = ctx
	return store
}
