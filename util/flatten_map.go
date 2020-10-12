package util

import (
	"strings"
)

// FlattenMap implements node.FlattenMap with key's usage record
type FlattenMap struct {
	keys        []string
	data        map[string]string
	used        map[string]bool
	originalKey map[string]string
}

// NewFlattenMap initialize a flatten map
func NewFlattenMap() *FlattenMap {
	return &FlattenMap{
		data:        make(map[string]string),
		used:        make(map[string]bool),
		originalKey: make(map[string]string),
	}
}

// Value implements node.FlattenMap with recording key's usage
func (m *FlattenMap) Value(key string) string {
	m.used[key] = true
	return m.data[key]
}

// ChildrenByPrefix implements node.FlattenMap
func (m *FlattenMap) ChildrenByPrefix(prefix string) (keys []string) {
	for _, key := range m.keys {
		if strings.HasPrefix(key, prefix) {
			if key = strings.SplitN(key[len(prefix):], CharDot, 2)[0]; key != "" {
				keys = append(keys, key)
			}
		}
	}
	return UniqueStrings(keys)
}

// Set format key to dot-case and set key to value, which also clear key's usage
func (m *FlattenMap) Set(originalKey, value string) {
	key := ToDotCase(originalKey)
	if _, set := m.data[key]; set {
		for i := range m.keys {
			if m.keys[i] == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)
				break
			}
		}
	}
	m.keys = append(m.keys, key)
	m.originalKey[key] = originalKey
	m.data[key] = value
	delete(m.used, key)
}

// Keys returns keys in order by when they were set
func (m *FlattenMap) Keys() []string {
	return m.keys
}

// Values returns values in order by when they were set
func (m *FlattenMap) Values() []string {
	values := make([]string, len(m.keys))
	for i := range m.keys {
		values[i] = m.data[m.keys[i]]
	}
	return values
}

// UnusedKeys returns unused keys after they were set
func (m *FlattenMap) UnusedKeys(escapedKeys []string) (keys []string) {
	escaped := StringsToSet(escapedKeys)
	for _, key := range m.keys {
		if !(m.used[key] || escaped[key]) {
			keys = append(keys, m.originalKey[key])
		}
	}
	return keys
}
