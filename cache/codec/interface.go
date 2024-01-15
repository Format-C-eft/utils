package codec

type Codec interface {
	Unmarshal(data []byte) (interface{}, error)
	Marshal(src interface{}) ([]byte, error)
}
