package helpers

import (
	"github.com/helloeave/json"
	"net/http"
)

type Message struct {
	Message interface{}
}

type Result struct {
	Code int
	Data interface{}
}

type Options struct {
	CountData int
}

type Err struct {
	Code  int
	Error Message
}

type Resp struct {
	RespObj http.ResponseWriter
}

func (w Resp) SendMessage(data Message) {
	w.RespObj.WriteHeader(http.StatusOK)

	er := Result{http.StatusOK, data}
	enc := json.NewEncoder(w.RespObj)
	enc.SetNilSafeCollection(true)

	if err := enc.Encode(er); err != nil {
		panic(err)
	}
}

func (w Resp) SendResponse(data interface{}) {
	w.RespObj.WriteHeader(http.StatusOK)

	res := Result{http.StatusOK, data}

	r, _ := json.MarshalSafeCollections(res)

	_, err := w.RespObj.Write(r)

	if err != nil {
		panic(err)
	}
}

func (w Resp) SendResponseMap(data map[string]interface{}) {
	w.RespObj.WriteHeader(http.StatusOK)

	res := Result{http.StatusOK, data}

	r, _ := json.MarshalSafeCollections(res)

	_, err := w.RespObj.Write(r)

	if err != nil {
		panic(err)
	}
}

func (w Resp) SendBadRequest(data Message) {
	w.RespObj.WriteHeader(http.StatusBadRequest)

	er := Err{http.StatusBadRequest, data}

	if err := json.NewEncoder(w.RespObj).Encode(er); err != nil {
		panic(err)
	}
}

func (w Resp) Send404() {
	w.RespObj.WriteHeader(http.StatusNotFound)

	data := Message{"ERROR: Не найден!"}
	er := Err{http.StatusNotFound, data}

	if err := json.NewEncoder(w.RespObj).Encode(er); err != nil {
		panic(err)
	}
}

func (w Resp) SendForbidden() {
	w.RespObj.WriteHeader(http.StatusForbidden)

	data := Message{"Forbidden: Ограничение доступа!"}
	er := Err{http.StatusForbidden, data}

	if err := json.NewEncoder(w.RespObj).Encode(er); err != nil {
		panic(err)
	}
}

func (w Resp) SendUnauthorized() {
	w.RespObj.WriteHeader(http.StatusUnauthorized)

	data := Message{"Expired token"}
	er := Err{http.StatusUnauthorized, data}

	if err := json.NewEncoder(w.RespObj).Encode(er); err != nil {
		panic(err)
	}
}

func (w Resp) SendNotFound(data Message) {
	w.RespObj.WriteHeader(http.StatusNotFound)

	er := Err{http.StatusNotFound, data}

	if err := json.NewEncoder(w.RespObj).Encode(er); err != nil {
		panic(err)
	}
}
