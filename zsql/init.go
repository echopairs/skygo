package zsql

import (
	"github.com/jmoiron/sqlx"

	"database/sql"
	"errors"
	"reflect"
)

var (
	ErrNoRows = sql.ErrNoRows
)

type (
	// VarTypeError indicates a variable type error when trying to populating a variable with DB result
	VarTypeError string
)

// Error returns the error message
func (s VarTypeError) Error() string {
	return "Invaild variable type: " + string(s)
}

type Rows struct {
	*sqlx.Rows
}

// dest must be the point of slice
func (db *DB) QueryAll(dest interface{}, sql string, args ...interface{}) error {
	// pointer
	vr := reflect.ValueOf(dest)
	if vr.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer not a value to QueryAll destination")
	}

	// not nil
	if vr.IsNil() {
		return errors.New("nil pointer passed to QueryAll destination")
	}

	// slice
	vr = reflect.Indirect(vr)
	if vr.Kind() != reflect.Slice {
		return errors.New("must pass Slice to QueryAll destination")
	}
	return db.Select(dest, sql, args...)
}

// dest must be point
func (db *DB) QueryOne(dest interface{}, sql string, args ...interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to QueryOne destination")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to QueryOne destination")
	}
	return db.Get(dest, sql, args...)
}

func (db *DB) QueryByName(arg interface{}, sql string) (*Rows, error) {
	rows, err := db.NamedQuery(sql, arg)
	if err != nil {
		return nil, err
	}
	return &Rows{rows}, nil
}

// Insert 插入數據，可以用來插入多行
// data must be : struct, slice-of-struct, slice-of-ptr2struct or ptr to these three
func (db *DB) Inserts(data interface{}, sql string) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.New("nil pointer found")
		}
		v = v.Elem()
	}

	var et reflect.Type
	var evs []reflect.Value // slice of structs
	if v.Kind() == reflect.Slice {
		et = v.Type().Elem()
		if et.Kind() == reflect.Ptr {
			et = et.Elem()
		}
		if et.Kind() != reflect.Struct {
			return VarTypeError("must be slice of struct or slice of ptr2struct")
		}
		for i := 0; i < v.Len(); i++ {
			ev := v.Index(i)
			if ev.Kind() == reflect.Ptr {
				ev = ev.Elem()
			}
			evs = append(evs, ev)
		}
	} else {
		// single
		et = v.Type()
		if et.Kind() != reflect.Struct {
			return VarTypeError("must be slice of struct or slice of ptr2struct")
		}
		evs = append(evs, v)
	}

	for _, ev := range evs {
		_, err := db.NamedExec(sql, ev.Interface())
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
