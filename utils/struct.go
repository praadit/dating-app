package utils

import (
	"encoding/json"
	"errors"
	"reflect"
)

func StructToMap(data any) (map[string]interface{}, error) {
	kind := reflect.TypeOf(data).Kind()
	if kind != reflect.Struct {
		return nil, errors.New("Data is not struct!")
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	structMap := map[string]interface{}{}
	err = json.Unmarshal(dataByte, &structMap)
	if err != nil {
		return nil, err
	}

	return structMap, err
}
func GetMapKeys[T comparable](maps map[T]interface{}) (slice []T) {
	for key := range maps {
		slice = append(slice, key)
	}
	return
}
func Contains[T comparable](slices []T, search T) bool {
	for _, key := range slices {
		if key == search {
			return true
		}
	}
	return false
}
