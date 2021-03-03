package model

import (
	"encoding/json"
	"fmt"
)

type ColdRecord struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"EndTime"`
	Hash      string `json:"hash"`
}

func CreateColdRecordKey(startTime string, endTime string) string {
	return MakeKey(startTime, endTime)
}

func (h *ColdRecord) GetSplitKey() []string {
	return []string{h.StartTime, h.Hash, h.EndTime}
}

func (h *ColdRecord) Serialize() ([]byte, error) {
	return json.Marshal(h)
}

func Deserialize(bytes []byte, hash *ColdRecord) error {
	err := json.Unmarshal(bytes, hash)

	if err != nil {
		return fmt.Errorf("Error deserializing commercial paper. %s", err.Error())
	}

	return nil
}
