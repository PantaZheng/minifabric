package store

import "marbles/model"

type ColdStoreInterface interface {
	AddColdRecord(*model.ColdRecord) error
	GetColdRecord(string, string) (*model.ColdRecord, error)
}

type coldStore struct {
	stateList StateListInterface
}

func (cs *coldStore) AddColdRecord(cr *model.ColdRecord) error {
	return cs.stateList.AddState(cr)
}

func (cs *coldStore) GetColdRecord(startTime string, endTime string) (*model.ColdRecord, error) {
	ch := new(model.ColdRecord)
	if err := cs.stateList.GetState(model.CreateColdRecordKey(startTime, endTime), ch); err != nil {
		return nil, err
	}
	return ch, nil
}

func newList(ctx model.TransactionContextInterface) *coldStore {
	stateList := new(StateList)
	stateList.Ctx = ctx

	stateList.Name = "cold record history"
	stateList.Deserialize = func(bytes []byte, state StateInterface) error {
		return model.Deserialize(bytes, state.(*model.ColdRecord))
	}

	list := new(coldStore)
	list.stateList = stateList
	return list
}
