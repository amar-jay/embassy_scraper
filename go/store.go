package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type Store[T Data] interface {
	// store to io.Writer
	Save() error
	Get(key string) (T, error)
	GetAll() ([]T, error)
	GetMap() (map[string]T, error)
	Fetch() error // refetch from store and reset store
	Set(key string, val T) Store[T]
	Del(key string) error
	Flush() error
	Close() error
	// GetFrom[T any](r io.Reader, data *T) error // get all data from reader
}

type MemStore[T Data] struct {
	data   map[string]T
	writer *bufio.Writer
	reader *bufio.Reader
	sync.RWMutex
}

type DataWrapper[T Data] struct {
	Data map[string]T `json:"data"`
}

type Data interface{}

// info about each embassy
type AccountInfo struct {
	Name       string   `json:"name"` // name of embassy
	Email      []string `json:"email"`
	Phone      []string `json:"phone"`
	Ambassador string   `json:"ambassador"`
	Address    string   `json:"address"` // location of embassy
}

var _ Data = (*AccountInfo)(nil)
var _ Store[AccountInfo] = (*MemStore[AccountInfo])(nil)

var (
	ErrKeyNotFound = errors.New("key not found")
)

func NewStore[T Data](file string) Store[T] {

	perm := fs.FileMode(0777)
	file = fmt.Sprintf("store/%s", strconv.Itoa(rand.Int())+"_"+file)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, perm)
	// f, err := os.CreateTemp("store", "*.json")
	if err != nil {
		panic(err)
	}
	// check if file write permissions are available if not, add them
	_, err = os.Stat(f.Name())
	if err != nil {
		if os.IsPermission(err) {
			err = os.Chmod(f.Name(), 0777)
			if err != nil {
				panic(err)
			}
		}
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	return &MemStore[T]{
		data:   make(map[string]T),
		writer: w,
		reader: bufio.NewReader(f),
	}
}

func (s *MemStore[T]) Get(key string) (T, error) {
	return s.data[key], nil
}
func (s *MemStore[T]) GetAll() ([]T, error) {
	var data []T
	s.RLock()
	defer s.RUnlock()
	for _, v := range s.data {
		data = append(data, v)
	}
	return data, nil
}

func (s *MemStore[T]) GetMap() (map[string]T, error) {
	s.RLock()
	defer s.RUnlock()
	return s.data, nil
}

func (s *MemStore[T]) Set(key string, val T) Store[T] {
	// log.Printf("Setting key: %s, val: %+v", key, val)
	s.Lock()
	s.data[key] = val
	s.Unlock()
	return Store[T](s)
}

func (s *MemStore[T]) Del(key string) error {
	if _, ok := s.data[key]; !ok {
		return ErrKeyNotFound
	}
	delete(s.data, key)
	return nil
}

func (s *MemStore[T]) Save() error {

	s.Lock()
	defer s.Unlock()
	datawrapped := DataWrapper[T]{
		Data: s.data,
	}

	d, err := json.Marshal(datawrapped)
	if err != nil {
		return err
	}
	d = append(d, '\n')

	_, err = s.writer.Write(d)

	if err != nil {
		return err
	}

	return nil
}

// get cached data from store
func (s *MemStore[T]) Fetch() error {
	// get all data from reader and unmarshal into store
	// read until end of file
	data, err := s.reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	var d DataWrapper[T]

	// check if it is valid json
	if !json.Valid(data) {
		return errors.New("invalid json")
	}
	err = json.Unmarshal(data, &d)
	fmt.Printf("%+v\n", d)
	s.Lock()
	s.data = d.Data
	s.Unlock()
	println(len(data), len(d.Data), len(s.data))
	if err != nil {
		panic(err)
		// return err
	}
	return nil

}

func (s *MemStore[T]) Flush() error {
	return s.writer.Flush()
}

// func (s *MemStore[T]) GetJSON(key string, store []T) error {
// 	// get all data from store and unmarshal into store
// 	data := make([]T, 0, 1024)
// 	for _, v := range s.data {
// 		data = append(data, v...)
// 	}
// 	err := json.Unmarshal(data, store)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (s *MemStore[T]) Close() error {
	err := s.writer.Flush()
	s.reader = nil
	s.writer = nil
	return err
}
