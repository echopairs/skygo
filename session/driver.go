package session

import (
	"sort"
	"sync"
)

var (
	drivers   = make(map[string]SessionDriver)
	driversMu sync.Mutex
)

// Store contains all data for one session process with specific id.
type SessionStore interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{})           // delete session value
	SessionID() string                // back current sessionID
}

// Driver contains global SessionStore methods and saved SessionStores.
// it can operate a SessionStore by its id
type SessionDriver interface {
	SessionCreate(sid string) (SessionStore, error)
	SessionDestroy(sid string) error
	SessionRead(sid string) (SessionStore, error)
	SessionGC(tickTime int64)
	SessionExist(sid string) bool
	SessionAll() []string
	SessionInit(cfg *ManagerConfig) error
}

// Register makes a session driver available by the provided name.
// if Register is called twice with the same name or if driver is nil.
// it panics.
func Register(name string, driver SessionDriver) {
	if driver == nil || name == "" {
		panic("session: Register driver is nil or name is nil")
	}
	driversMu.Lock()
	defer driversMu.Unlock()
	if _, dup := drivers[name]; dup {
		panic("session: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Drivers returns a sorted list of the names of the registered drivers
func Drivers() []string {
	driversMu.Lock()
	defer driversMu.Lock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
