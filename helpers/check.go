package helpers

import (
	"errors"
	"reflect"
)

func HasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct || rv.Kind() != reflect.Map {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

func HasMethod(v interface{}, methodName string) bool {
	st := reflect.TypeOf(v)
	_, ok := st.MethodByName(methodName)
	return ok
}

func CallMethod(v interface{}, methodName string) (interface{}, error) {
	if HasMethod(v, methodName) {
		l := reflect.ValueOf(&v).MethodByName(methodName)
		e := l.Call([]reflect.Value{})
		return e, nil
	} else {
		return nil, errors.New("Method does not exists")
	}
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
