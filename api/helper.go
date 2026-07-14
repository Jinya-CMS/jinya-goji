package api

import (
	"encoding/json"
	"net/http"
)

func ParseBody[T any](r *http.Request) (*T, error) {
	decoder := json.NewDecoder(r.Body)
	var body T
	err := decoder.Decode(&body)

	return &body, err
}

type Body struct {
	data map[string]any
}

func BodyFromRequest(r *http.Request) (*Body, error) {
	decoder := json.NewDecoder(r.Body)
	var body map[string]any
	err := decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return new(Body{data: body}), nil
}

func GetFromBody[T any](body *Body, key string) (*T, bool) {
	res, ok := body.data[key]
	if ok {
		val, ok := res.(T)
		return &val, ok
	}

	return nil, false
}

func GetIntFromBody(body *Body, key string) (int, bool) {
	res, ok := body.data[key]
	if ok {
		val, ok := res.(float64)
		return int(val), ok
	}

	return -1, false
}

func GetArrayFromBody[T any](body *Body, key string) ([]T, bool) {
	res, ok := body.data[key]
	if ok {
		data, ok := res.([]interface{})
		if !ok {
			return nil, false
		}
		val := make([]T, len(data))
		for i, v := range data {
			val[i], ok = v.(T)
			if !ok {
				return nil, false
			}
		}

		return val, ok
	}

	return nil, false
}
