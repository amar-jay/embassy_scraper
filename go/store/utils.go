package store

import (
	"reflect"
)

// get keys from struct type T
func GetKeys[T Data](s T) []string {
	var keys []string
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		keys = append(keys, v.Type().Field(i).Name)
	}
	return keys
}

// get keys with prefix from struct type T
func GetKeysWithPrefix[T Data](s T, prefix string) []string {
	var keys []string
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		keys = append(keys, prefix+v.Type().Field(i).Name)
	}
	return keys
}
