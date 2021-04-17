package store

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const collectionName = "collectionPrivate"
const indexStart = 10000

//var emptyValue []byte

//type RetentionPolicy struct {
//	StartTime timestamp.Timestamp `json:"start_time"`
//	Interval  int64               `json:"interval"`
//}

type Record struct {
	Timestamp int    `json:"timestamp"`
	DeviceID  string `json:"device_id"`
	Data      []byte `json:"data"`
}

type TstData struct {
	Data string `json:"data"`
}

func (r *Record) makePrimaryKey() string {
	return string(rune(indexStart + r.Timestamp))
}

type HotStoreInterface interface {
	getTransient(record *Record) error
	AddPubRecord(record *Record) error
	AddPvtRecord(record *Record) error
	AddPubRecordTst(record *Record) error
	AddPvtRecordTst(record *Record) error
	GetPubRecord(record *Record) error
	GetPvtRecord(record *Record) error
	AddPubPvtTst(record *Record) error
	HSSQ(number int) error
	PSRQ(number int) error
	HSMQ(number int) error
	GS() error
	GPD() error
	GPDH() error
	//GetRecordHash(timestamp, id string) (string, error)
	//GetRecordsByRange(startKey, endKey string) ([]Record, error)
}

type hotStore struct {
	Ctx contractapi.TransactionContextInterface
}

func (hs *hotStore) getTransient(record *Record) error {
	transMap, err := hs.Ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	transientJSON, ok := transMap["time"]
	if !ok {
		return fmt.Errorf("hot_record not found in the transient map")
	}

	var tstData TstData
	err = json.Unmarshal(transientJSON, &tstData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}
	record.Data = []byte(tstData.Data)
	return err
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

func (hs *hotStore) AddPubRecordTst(record *Record) error {
	err := hs.getTransient(record)
	if err != nil {
		return fmt.Errorf("AddPubRecordTst-failed to get transient: %s", err.Error())
	}

	err = hs.AddPubRecord(record)
	if err != nil {
		return fmt.Errorf(" AddPubRecordTst- %s", err.Error())
	}

	return err
}

func (hs *hotStore) AddPvtRecordTst(record *Record) error {
	err := hs.getTransient(record)
	if err != nil {
		return fmt.Errorf("AddPvtRecordTst-failed to get transient: %s", err.Error())
	}

	err = hs.AddPvtRecord(record)
	if err != nil {
		return fmt.Errorf(" AddPvtRecordTst- %s", err.Error())
	}

	return err
}

func (hs *hotStore) AddPubPvtTst(record *Record) error {
	err := hs.getTransient(record)
	if err != nil {
		return fmt.Errorf("AddPvtRecordTst-failed to get transient: %s", err.Error())
	}

	err = hs.AddPvtRecord(record)
	if err != nil {
		return fmt.Errorf(" AddPubRecordTst- %s", err.Error())
	}
	hash, err := hs.Ctx.GetStub().GetPrivateDataHash(collectionName, string(rune(record.Timestamp)))
	if err != nil {
		return fmt.Errorf(" get data_hash- %s", err.Error())
	}
	record.Data = hash
	err = hs.AddPubRecord(record)
	if err != nil {
		return fmt.Errorf(" AddPvtRecordTst- %s", err.Error())
	}

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

func (hs *hotStore) GS() error {
	_, err := hs.Ctx.GetStub().GetState(string(rune(indexStart + 1)))
	if err != nil {
		return fmt.Errorf("getPrivateDataHash error %s", err.Error())
	}
	return nil
}

func (hs *hotStore) GPD() error {
	_, err := hs.Ctx.GetStub().GetPrivateData(collectionName, string(rune(indexStart+1)))
	if err != nil {
		return fmt.Errorf("getPrivateDataHash error %s", err.Error())
	}
	return nil
}

func (hs *hotStore) GPDH() error {
	_, err := hs.Ctx.GetStub().GetPrivateDataHash(collectionName, string(rune(indexStart+1)))
	if err != nil {
		return fmt.Errorf("getPrivateDataHash error %s", err.Error())
	}
	return nil
}

func (hs *hotStore) HSSQ(number int) error {
	for i := 0; i < number; i++ {
		_, err := hs.Ctx.GetStub().GetPrivateDataHash(collectionName, string(rune(indexStart+i)))
		if err != nil {
			return fmt.Errorf("getPrivateDataHash error %s", err.Error())
		}
	}
	return nil
}

func (hs *hotStore) HSMQ(number int) error {
	_, err := hs.Ctx.GetStub().GetPrivateDataByRange(
		collectionName, string(rune(indexStart)), string(rune(indexStart+number)))
	if err != nil {
		return fmt.Errorf("GetPrivateDataByRange error %s", err.Error())
	}
	return nil
}

func (hs *hotStore) PSRQ(number int) error {
	_, err := hs.Ctx.GetStub().GetStateByRange(
		string(rune(indexStart)), string(rune(indexStart+number)))
	if err != nil {
		return fmt.Errorf("GetStateByRange error %s", err.Error())
	}
	return nil
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
