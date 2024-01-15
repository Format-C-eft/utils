package codec

import (
	"encoding/json"
	"fmt"
)

type Json struct{}

func (Json) Marshal(in interface{}) ([]byte, error) {
	marshaller, ok := in.(json.Marshaler)
	if ok {
		val, err := marshaller.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshaler.codec.json err: %w", err)
		}
		return val, nil
	}

	val, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("marshal.codec.json err: %w", err)
	}
	return val, nil
}

func (Json) Unmarshal(p []byte, dst interface{}) error {
	unmarshaler, ok := dst.(json.Unmarshaler)
	if ok {
		err := unmarshaler.UnmarshalJSON(p)
		if err != nil {
			return fmt.Errorf("unmarshaler.codec.json err: %w", err)
		}

		return nil
	}

	err := json.Unmarshal(p, dst)
	if err != nil {
		return fmt.Errorf("unmarshal.codec.json err: %w", err)
	}
	return nil
}
