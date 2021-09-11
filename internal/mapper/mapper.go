package mapper

// Mapper is a map with some useful methods
type Mapper map[string]interface{}

// AddValue adds a value to the Mapper
func (m Mapper) AddValue(key string, value interface{}) {
	m[key] = value
}

// GetValue get the value in the map for a given key, if it exists.
func (m Mapper) GetValue(key string) (value interface{}, exists bool) {
	if value, exists := m[key]; exists {
		return value, true
	}

	return nil, false
}

// Each executes a given function iterating over all key/value pairs in the Mapper
func (m Mapper) Each(fn func(key string, value interface{})) {
	for k, v := range m {
		fn(k, v)
	}
}
