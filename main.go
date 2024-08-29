package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Storer[K comparable, V any] interface {
	Get(K) (V, error)
	Put(K, V) error
	Update(K, V) error
	Delete(K) (V, error)
}

type KVStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		data: make(map[K]V),
	}
}

// this function is NOT concurrent-safe
func (s *KVStore[K, V]) Has(key K) bool {
	_, ok := s.data[key]
	return ok
}

func (s *KVStore[K, V]) Update(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Has(key) {
		return fmt.Errorf("the key (%v) does not exist", key)
	}

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("value for key (%v) does not exist", key)
	}

	return value, nil

}

func (s *KVStore[K, V]) Delete(key K) (V, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("the key (%v) does not exist", key)
	}

	delete(s.data, key)

	return value, nil

}

type User struct {
	ID        int
	FirstName string
	Age       int
	Gender    string
}

type Server struct {
	Storage       Storer[int, *User]
	listenAddress string
}

func NewServer(listenAddress string) *Server {
	return &Server{Storage: NewKVStore[int, *User](), listenAddress: listenAddress}
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("handlePut function checking in!"))
}

func (s *Server) Start() {
	fmt.Printf("HTTP Server now running on port (%s)", s.listenAddress)

	http.HandleFunc("/put", s.handlePut)

	log.Fatal(http.ListenAndServe(s.listenAddress, nil))
}

func main() {
	s := NewServer(":3000")

	s.Start()
	// store := NewKVStore[string, string]()

	// if err := store.Put("foo", "bar"); err != nil {
	// 	log.Fatal(err)
	// }

	// value, err := store.Get("foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(value)

	// if err := store.Update("foo", "oof"); err != nil {
	// 	log.Fatal(err)
	// }

	// value, err = store.Get("foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(value)

	// if value, err = store.Delete("foo"); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(value)
}
