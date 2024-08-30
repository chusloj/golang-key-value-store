package main

import "testing"

func TestPut(t *testing.T) {
	s := NewServer(":3000")

	key := "foo"
	value := "bar"

	s.Storage.Put(key, value)
	retrieved_value, err := s.Storage.Get(key)
	if err != nil {
		t.Errorf("key (%s) doesn't exist", key)
	}

	if retrieved_value != value {
		t.Errorf("retreived value %s does not equal intended value %s", retrieved_value, value)
	}
}
