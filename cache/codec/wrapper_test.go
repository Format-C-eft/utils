package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapperJSON(t *testing.T) {
	t.Parallel()
	type testType struct {
		TestA string `json:"test_a"`
		TestB string `json:"test_b"`
	}

	wr := NewWrapped(
		func() interface{} { return &testType{} },
		Json{},
	)

	testStruct := testType{
		TestA: "testA",
		TestB: "testB",
	}
	testByte := []byte("{\"test_a\":\"testA\",\"test_b\":\"testB\"}")

	t.Run("Marshal", func(t *testing.T) {
		t.Parallel()
		data, err := wr.Marshal(testStruct)
		assert.Nil(t, err)
		assert.Equal(t, data, testByte)
	})

	t.Run("Unmarshal", func(t *testing.T) {
		t.Parallel()
		data, err := wr.Unmarshal(testByte)
		assert.Nil(t, err)
		assert.Equal(t, data, &testStruct)
	})

}
