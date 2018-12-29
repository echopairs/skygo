package zsql

import (
	"testing"
)

var sqlAddress = SqlAddress{
	"mysql",
	"127.0.0.1:3306",
	1,
	1,
	"root",
	"admin",
	"pairs",
	1,
}

type TestTable struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func TestConnect(t *testing.T) {
	_, err := Connect(&sqlAddress)
	if err != nil {
		t.Error(err)
		return
	}
}
