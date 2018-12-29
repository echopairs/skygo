package zsql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"fmt"
	"time"
)

var (
	drivers []string
)

type SqlAddress struct {
	DriverName   string `yaml:"driver_name"`
	Address      string `yaml:"address"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	UserName     string `yaml:"user_name"`
	Password     string `yaml:"password"`
	DbName       string `yaml:"db_name"`
	DialTimeout  int    `yaml:"dial_timeout"`
}

type DB struct {
	*sqlx.DB
	*SqlAddress
}

func init() {
	drivers = append(drivers, "mysql")
}

// Connect to a database and verify with a ping
func Connect(addr *SqlAddress) (*DB, error) {
	db, err := Open(addr)
	if err != nil {
		return nil, err
	}

	// verify
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// connect to a database return db pool
func Open(addr *SqlAddress) (*DB, error) {
	if err := checkSqlAddress(addr); err != nil {
		return nil, err
	}
	dataSourceName := addr.UserName + ":" + addr.Password + "@tcp(" + addr.Address + ")/" + addr.DbName
	db, err := sqlx.Open(addr.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(addr.MaxOpenConns)
	db.SetMaxIdleConns(addr.MaxIdleConns)
	db.SetConnMaxLifetime(time.Second * 60)
	return &DB{db, addr}, nil
}

func checkSqlAddress(addr *SqlAddress) error {
	// 1. check addr
	if addr == nil {
		return fmt.Errorf("zsql config is nil, please check")
	}
	if addr.Address == "" ||
		addr.DbName == "" ||
		addr.UserName == "" ||
		addr.Password == "" {
		return fmt.Errorf("sql addr cfg %v is error", addr)
	}
	// 2. check driver
	if exist := checkDriver(addr.DriverName); !exist {
		return fmt.Errorf("not suppurt %s driver", addr.DriverName)
	}
	return nil
}

func checkDriver(driverName string) bool {
	exist := false
	for _, d := range drivers {
		if d == driverName {
			exist = true
			break
		}
	}
	return exist
}
