package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
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

type Server struct {
	Storage       Storer[string, string]
	listenAddress string
}

func NewServer(listenAddress string) *Server {
	return &Server{Storage: NewKVStore[string, string](), listenAddress: listenAddress}
}

func (s *Server) handlePut(c echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")

	s.Storage.Put(key, value)

	return c.JSON(http.StatusOK, map[string]string{"msg": "ok"})
}

func (s *Server) handleGet(c echo.Context) error {
	key := c.Param("key")

	value, err := s.Storage.Get(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"value": value})
}

func (s *Server) Start() {
	fmt.Printf("HTTP Server now running on port (%s)", s.listenAddress)

	e := echo.New()

	// labeled as a GET route so that testing can be done in the browser
	e.GET("/put/:key/:value", s.handlePut)
	e.GET("/get/:key", s.handleGet)

	e.Start(s.listenAddress)

	// http.HandleFunc("/put", s.handlePut)
	// log.Fatal(http.ListenAndServe(s.listenAddress, nil))
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
