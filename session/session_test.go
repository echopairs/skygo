package session

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

var (
	manager *Manager
	cfg     = &ManagerConfig{
		CookieName:     "cloudId",
		GcLifeTime:     1,
		MaxLifeTime:    3,
		DriverName:     "memory",
		CookieLifeTime: 60 * 1,
	}
)

func TestNewManager(t *testing.T) {
	manager, err := NewManager(cfg)
	if err != nil {
		t.Error(err)
	}
	if manager == nil {
		t.Error("create manager error")
	}
}

type headerOnlyResponseWriter http.Header

func (ho headerOnlyResponseWriter) Header() http.Header {
	return http.Header(ho)
}

func (ho headerOnlyResponseWriter) Write([]byte) (int, error) {
	panic("NOIMPL")
}

func (ho headerOnlyResponseWriter) WriteHeader(int) {
	panic("NOIMPL")
}

func TestManager_SessionCreate(t *testing.T) {
	mgr, err := NewManager(cfg)
	if err != nil {
		t.Error(err)
	}
	manager = mgr
	//var w http.ResponseWriter
	r, err := http.NewRequest("", "//127.0.0.1:9090/", nil)
	if err != nil {
		t.Errorf("newRequest failed for %v", err.Error())
	}
	w := &headerOnlyResponseWriter{}
	session, err := manager.SessionCreate(w, r)
	if err != nil {
		t.Error(err)
	}

	s, err := manager.SessionRead(session.SessionID())
	if err != nil {
		t.Error("session read ok ", err)
	}
	time.Sleep(time.Second * 4)
	s, err = manager.SessionRead(session.SessionID())
	if err == nil {
		// GC
		t.Error("error gc is not working")
	}
	fmt.Printf(" %v", s)
}
