package database

import (
	"reflect"

	"github.com/DerKnerd/gorp"
)

func Get[R any](keys ...any) (*R, error) {
	var r R
	res, err := dbMap.Get(&r, keys...)
	if res != nil {
		return res.(*R), err
	}

	return nil, err
}

func SelectOne[R any](query string, keys ...any) (R, error) {
	var r R
	err := dbMap.SelectOne(&r, query, keys...)

	return r, err
}

func Select[R any](query string, keys ...any) ([]R, error) {
	r := new(R)
	res, err := dbMap.Select(&r, query, keys...)
	if err != nil {
		return nil, err
	}

	result := make([]R, len(res))
	for idx, item := range res {
		elem := item.(*R)
		result[idx] = *elem
	}

	return result, nil
}

func toPtr[R any](value any) any {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		return value
	} else {
		elem := value.(R)
		return &elem
	}
}

func AddTableWithName[R any](name string) *gorp.TableMap {
	var r R
	return GetDbMap().AddTableWithName(r, name)
}
