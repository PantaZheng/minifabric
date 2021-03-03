package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var EmptyValue = []byte{0x00}

type StoreInterface interface {
	makeIndexKey()
	PutPublicData(pub PublicData)
	PutPrivateData(pvt PrivateData)
	PutPublicDataWithIndex(pub PublicData)
	PutPrivateDataWithIndex(pvt PrivateData)
	GetPublicData(pub *PublicData)
	GetPrivateData(pvt *PrivateData)
	GetPublicDataByRange(startKey, endKey string)
	GetPrivateDataByRange(col, startKey, endKey string)
	GetPublicDataHistory(pub *PublicData)
}

type Index struct {
	Name  string
	Attrs []string
}

type PublicData struct {
	Key   string
	Value []byte
	Index Index
}

type PrivateData struct {
	*PublicData
	ColName string
	Hash    []byte
}

type Store struct {
	Ctx contractapi.TransactionContextInterface
}

func (s *Store) makeIndexKey(index Index) error {
	indexKey, err := s.Ctx.GetStub().CreateCompositeKey(index.Name, index.Attrs)
	if err != nil {
		return fmt.Errorf("failed to enable index: %s", err.Error())
	}

	if err := s.Ctx.GetStub().PutState(indexKey, EmptyValue); err != nil {
		return fmt.Errorf("failed to put index: %s", err.Error())
	}
	return nil
}

func (s *Store) PutPublicData(pub PublicData) error {
	if err := s.Ctx.GetStub().PutState(pub.Key, pub.Value); err != nil {
		return fmt.Errorf("failed to put public data: %s", err.Error())
	}
	return nil
}

func (s *Store) PutPrivateData(pvt PrivateData) error {
	var err error
	if err = s.Ctx.GetStub().PutPrivateData(pvt.ColName, pvt.Key, pvt.Value); err != nil {
		return fmt.Errorf("failed to put private data: %s", err.Error())
	}
	if pvt.Hash, err = s.Ctx.GetStub().GetPrivateDataHash(pvt.ColName, pvt.Key); err != nil {
		return fmt.Errorf("failed to get private data: %s", err.Error())
	}
	return nil
}

func (s *Store) PutPublicDataWithIndex(pub PublicData) error {
	if err := s.PutPublicData(pub); err != nil {
		return err
	}
	if err := s.makeIndexKey(pub.Index); err != nil {
		return err
	}
	return nil
}

func (s *Store) PutPrivateDataWithIndex(pvt PrivateData) error {
	if err := s.PutPrivateData(pvt); err != nil {
		return err
	}
	if err := s.makeIndexKey(pvt.Index); err != nil {
		return err
	}
	return nil
}

func (s *Store) GetPublicData(pub *PublicData) error {
	var err error
	pub.Value, err = s.Ctx.GetStub().GetState(pub.Key)
	return err
}

func (s *Store) GetPrivateData(pvt *PrivateData) error {
	var err error
	if pvt.Value, err = s.Ctx.GetStub().GetPrivateData(pvt.ColName, pvt.Key); err != nil {
		return err
	}
	if pvt.Hash, err = s.Ctx.GetStub().GetPrivateDataHash(pvt.ColName, pvt.Key); err != nil {
		return err
	}
	return err
}

func (s *Store) GetPublicDataByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	return s.Ctx.GetStub().GetStateByRange(startKey, endKey)
}

func (s *Store) GetPrivateDataByRange(col, startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	return s.Ctx.GetStub().GetPrivateDataByRange(col, startKey, endKey)
}

func (s *Store) GetPublicDataHistory(pub *PublicData) (shim.HistoryQueryIteratorInterface, error) {
	return s.Ctx.GetStub().GetHistoryForKey(pub.Key)
}
