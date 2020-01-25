package helpers

import (
	"errors"
	"fmt"
	"reflect"
)

type ResMap map[string]interface{}

func IndexOf(value string, arr map[int]string) int {
	for index, element := range arr {
		if element == value {
			return index
		}
	}

	return -1
}

func GetStructName(p interface{}) string {
	return reflect.TypeOf(p).Name()
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(obj interface{}, m map[string]interface{}) (interface{}, error) {
	err := new(error)
	for k, v := range m {
		err := SetField(obj, k, v)
		if err != nil {
			return obj, err
		}
	}
	return obj, *err
}

func UnionMaps(m1, m2 map[string]interface{}) map[string]interface{} {
	for ia, va := range m1 {
		if _, ok := m2[ia]; ok {
			continue
		}
		m2[ia] = va

	}
	return m2
}
