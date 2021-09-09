package singertap

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/alecthomas/jsonschema"
)

type SingerTap struct {
	sync.Mutex

	currentRecordType string // need to track when the record type changes
}

func New() *SingerTap {
	s := SingerTap{}
	return &s
}

func (s *SingerTap) WriteRecord(recordType string, record Record) error {
	out, err := json.Marshal(record)
	if err != nil {
		return err
	}

	r := SingerRecord{
		SingerMessage: SingerMessage{
			Type:   "RECORD",
			Stream: recordType,
		},
		Record: out,
	}

	out, err = json.Marshal(r)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	if recordType != s.currentRecordType {
		err := s.writeSchema(recordType, record)
		if err != nil {
			return err
		}
		s.currentRecordType = recordType
	}

	_, err = fmt.Fprint(os.Stdout, string(out))
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(os.Stdout, "\n")
	if err != nil {
		return err
	}

	return nil
}

func (s *SingerTap) writeSchema(recordType string, record Record) error {

	reflector := jsonschema.Reflector{
		ExpandedStruct: true,
	}

	schema := reflector.Reflect(record)
	out, err := schema.MarshalJSON()
	if err != nil {
		return err
	}

	r := SingerSchema{
		SingerMessage: SingerMessage{
			Type:   "SCHEMA",
			Stream: recordType,
		},
		KeyProperties: record.KeyProperties(),
		Schema:        out,
	}

	out, err = json.Marshal(r)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, string(out))
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(os.Stdout, "\n")
	if err != nil {
		return err
	}

	return nil
}

type SingerMessage struct {
	Type   string `json:"type"`
	Stream string `json:"stream"`
}

type SingerRecord struct {
	SingerMessage
	Record json.RawMessage `json:"record"`
}

type SingerSchema struct {
	SingerMessage
	KeyProperties []string        `json:"key_properties"`
	Schema        json.RawMessage `json:"schema"`
}

type Record interface {
	KeyProperties() []string
}
