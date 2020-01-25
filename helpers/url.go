package helpers

import (
	"net/url"
	"reflect"
	"strings"
)

func prepareKey(key string, value string) map[string]interface{} {
	var m = make(map[string]interface{})
	var subm = make(map[string]interface{})
	var v = []string{}
	if i := strings.Index(key, "["); i >= 0 {
		subkey := key[i+1:]
		if i := strings.Index(subkey, "]"); i >= 0 {
			subkey = subkey[0:i]
		}
		key = key[:i]

		if i := strings.Index(value, "["); i >= 0 {
			subm[subkey] = prepareValue(value, v)
		} else {
			subm[subkey] = value
		}

		m[key] = subm

	} else {
		if i := strings.Index(value, "["); i >= 0 {
			m[key] = prepareValue(value, v)
		} else {
			m[key] = value
		}
	}

	return m
}

func prepareValue(value string, v []string) []string {
	value = value[1:]
	if i := strings.Index(value, "]"); i >= 0 {
		value = value[:i]
	}
	if i := strings.Index(value, ","); i >= 0 {
		v = append(v, value[:i])
		value = value[i:]
		value = strings.TrimSpace(value)
		v = prepareValue(value, v)
	} else {
		value = strings.TrimSpace(value)
		v = append(v, value)
	}

	return v
}

func ParseQuery(query string) map[string]interface{} {
	var m = make(map[string]interface{})
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}

		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			continue
		}

		if value == "" {
			continue
		}

		ar := prepareKey(key, value)
		for k, v := range ar {
			if _, ok := m[k]; ok && reflect.TypeOf(v).String() == reflect.TypeOf(m).String() {
				v4, ok := m[k].(map[string]interface{})
				v3, ok := v.(map[string]interface{})
				if !ok {
					// Can't assert, handle error.
				}
				for k2, v2 := range v3 {
					v4[k2] = v2
				}
			} else {
				m[k] = v
			}
		}

	}
	return m
}
