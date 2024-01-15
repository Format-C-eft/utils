package codec

import (
	"fmt"
)

type RawCodec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type NewObjectFN func() interface{}

type wrapped struct {
	codec  RawCodec
	newObj NewObjectFN
}

func (c wrapped) Unmarshal(payload []byte) (interface{}, error) {
	obj := c.newObj()
	err := c.codec.Unmarshal(payload, obj)
	if err != nil {
		return nil, fmt.Errorf("unmarshal.codec.wrapped err: %w", err)
	}

	return obj, nil
}

func (c wrapped) Marshal(src interface{}) ([]byte, error) {
	val, err := c.codec.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("marshal.codec.wrapped err: %w", err)
	}

	return val, nil
}

func NewWrapped(newFN NewObjectFN, codec RawCodec) Codec {
	return wrapped{
		codec:  codec,
		newObj: newFN,
	}
}
