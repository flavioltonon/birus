package shingling

import (
	"bytes"
	"encoding/gob"

	"github.com/flavioltonon/birus/internal/mapper"
)

// ShinglesMapper is a mapper for Shingles by their hashes
type ShinglesMapper struct {
	m mapper.Mapper
}

// NewShinglesMapper creates a new ShinglesMapper
func NewShinglesMapper() *ShinglesMapper {
	return &ShinglesMapper{
		m: make(mapper.Mapper),
	}
}

// AddValue adds a new Shingle to the ShinglesMapper
func (sm *ShinglesMapper) AddValue(key string, value *Shingle) {
	sm.m.AddValue(key, value)
}

// GetValue returns a Shingle identified by a given key in the ShinglesMapper
func (sm *ShinglesMapper) GetValue(key string) (*Shingle, bool) {
	if v, exists := sm.m.GetValue(key); exists {
		return v.(*Shingle), true
	}

	return nil, false
}

// Each executes a given function iterating over all key/value pairs in the ShinglesMapper
func (sm *ShinglesMapper) Each(fn func(key string, value *Shingle)) {
	sm.m.Each(func(k string, v interface{}) { fn(k, v.(*Shingle)) })
}

// Length returns the length of the ShinglesMapper
func (sm *ShinglesMapper) Length() int {
	return len(sm.m)
}

type GobShinglesMapper struct {
	M mapper.Mapper
}

func (m *ShinglesMapper) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(&GobShinglesMapper{M: m.m}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (m *ShinglesMapper) GobDecode(b []byte) error {
	var (
		buffer = bytes.NewBuffer(b)
		reader GobShinglesMapper
	)

	if err := gob.NewDecoder(buffer).Decode(&reader); err != nil {
		return err
	}

	m.m = reader.M
	return nil
}
