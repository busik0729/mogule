package helpers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const ACCESS_TOKEN, REFRESH_TOKEN, FORCE, DEVICE = "X-AT", "X-RT", "X-Force", "X-Device"

type Headers map[string]interface{}

func getHeaders(r *http.Request) Headers {
	headers := make(Headers)

	for k, v := range r.Header {
		headers[strings.ToLower(k)] = string(v[0])
	}

	return headers
}

func GetURLHeaderByKey(r *http.Request, key string) string {
	headers := getHeaders(r)
	key = strings.ToLower(key)

	if value, ok := headers[key]; ok {
		return value.(string)
	}

	return ""
}

func GetAccessTokenHeader(r *http.Request) string {
	authToken := GetURLHeaderByKey(r, ACCESS_TOKEN)
	authToken = strings.TrimSpace(authToken)

	return authToken
}

func GetAccessTokenByGetParam(r *http.Request) string {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		panic(err)
	}

	m := ParseQuery(u.RawQuery)
	if value, ok := m[ACCESS_TOKEN]; ok {
		return strings.TrimSpace(value.(string))
	}

	return ""
}

func GetRefreshTokenHeader(r *http.Request) string {
	authToken := GetURLHeaderByKey(r, REFRESH_TOKEN)
	authToken = strings.TrimSpace(authToken)

	return authToken
}

func GetDeviceHeader(r *http.Request) string {
	device := GetURLHeaderByKey(r, DEVICE)
	device = strings.TrimSpace(device)

	return device
}

func GetDeviceValues(r *http.Request) ([]string, error) {
	device := GetDeviceHeader(r)
	if len(device) < 1 {
		return []string{}, errors.New("Headers error")
	}
	dev := strings.Split(device, "!")

	return dev, nil
}

func GetForceHeader(r *http.Request) string {
	force := GetURLHeaderByKey(r, FORCE)
	force = strings.TrimSpace(force)

	return force
}
