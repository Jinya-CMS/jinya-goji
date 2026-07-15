package database

import (
	"log/slog"

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

func SelectOne[R any](query string, keys ...any) (*R, error) {
	var r R
	err := dbMap.SelectOne(&r, query, keys...)
	if err != nil {
		return nil, err
	}

	return &r, nil
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

func Insert(value any) error {
	err := dbMap.Insert(value)
	if err != nil {
		slog.Error("Database insert failed", "error", err, "value", value)
	}
	return err
}

func Update(value any) error {
	_, err := dbMap.Update(value)
	if err != nil {
		slog.Error("Database update failed", "error", err, "value", value)
	}
	return err
}

func Delete(value any) error {
	_, err := dbMap.Delete(value)
	if err != nil {
		slog.Error("Database delete failed", "error", err, "value", value)
	}
	return err
}

func AddTableWithName[R any](name string) *gorp.TableMap {
	var r R
	return GetDbMap().AddTableWithName(r, name)
}
