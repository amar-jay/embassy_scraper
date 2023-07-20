package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Store interface {
	// store to io.Writer
	StoreTo(w io.Writer) error
	Get(key string) ([]byte, error)
	Reload() error // refetch from store and reset store
	Set(key string, val []byte) error
	Del(key string) error
	// GetFrom[T any](r io.Reader, data *T) error // get all data from reader
}

type MemStore struct {
	data map[string]Data
	writer io.Writer
	reader io.Reader
}

type Data interface{}

// info about each embassy
type AccountInfo struct {
	Email   string `json:"email"`
	Phone1  string `json:"phone1"`
	Phone2  string `json:"phone2"`
	Phone3  string `json:"phone3"`
	Address string `json:"address"` // location of embassy
	Fax     string `json:"fax"`
}

var _ Data = (*AccountInfo)(nil)

var (
	ErrKeyNotFound = errors.New("key not found")
)

func NewStore[T Data](file string, dataType T) Store {

	// create a new reader
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	w := bufio.NewWriter(f)
	return &MemStore{
		data:   make(map[string]T),
		writer: w,
		reader: r,
	}
}

func (s *MemStore) Get(key string) ([]byte, error) {
	return s.data[key], nil
}

func (s *MemStore) Set(key string, val []byte) error {
	s.data[key] = val
	return nil
}

func (s *MemStore) Del(key string) error {
	if _, ok := s.data[key]; !ok {
		return ErrKeyNotFound
	}
	delete(s.data, key)
	return nil
}

func (s *MemStore) StoreTo(f io.Writer) error {

	// create a new reader
	reader := bufio.NewWriter(f)
	data := make([]byte, 0, 1024)

	d, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	copy(data, d[:])
	// write to file
	_, err = reader.WriteString(string(data))
	if err != nil {
		return err
	}

	return nil
}

func (s *MemStore) GetJSON(key string, store *[]byte) error {
	// get all data from store and unmarshal into store
	data := make([]byte, 0, 1024)
	for _, v := range s.data {
		data = append(data, v...)
	}
	err := json.Unmarshal(data, store)
	if err != nil {
		return err
	}
	return nil
}
