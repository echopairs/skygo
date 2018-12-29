package session

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/satori/go.uuid"
)

type ManagerConfig struct {
	CookieName     string	`yaml:"cookie_name"`
	GcLifeTime     int64	`yaml:"gclife_time"`
	MaxLifeTime    int64	`yaml:"maxlife_time"`
	DriverName     string	`yaml:"driver_name"`
	CookieLifeTime int		`yaml:"cookielife_time"`
}

type Manager struct {
	*ManagerConfig
	SessionDriver
}

func NewManager(cfg *ManagerConfig) (*Manager, error) {
	driver, ok := drivers[cfg.DriverName]
	if !ok {
		return nil, fmt.Errorf("session: unknown driver %q (forgotten import?)", cfg.DriverName)
	}
	if err := driver.SessionInit(cfg); err != nil {
		return nil, err
	}
	return &Manager{cfg, driver}, nil
}

func (manager *Manager) SessionCreate(w http.ResponseWriter, r *http.Request) (session SessionStore, err error) {
	sid, err := sessionId()
	if err != nil {
		return
	}
	if sid != "" && manager.SessionExist(sid) {
		return manager.SessionRead(sid)
	}

	// create
	session, err = manager.SessionDriver.SessionCreate(sid)
	if err != nil {
		fmt.Printf("create session failed because of %s ", err.Error())
		return
	}

	//fmt.Printf("create session sid is %s ", sid)

	cookie := &http.Cookie{
		Name:     manager.CookieName,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		Domain:   r.URL.Host,
		MaxAge:   manager.CookieLifeTime,
		Expires:  time.Now().Add(time.Duration(manager.CookieLifeTime) * time.Second),
	}
	http.SetCookie(w, cookie)
	return
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.CookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	// decoder
	sid, _ := url.QueryUnescape(cookie.Value)
	manager.SessionDriver.SessionDestroy(sid)
	expiration := time.Now()
	cookie = &http.Cookie{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  expiration,
	}
	http.SetCookie(w, cookie)
}

func sessionId() (string, error) {
	sid := uuid.NewV4()
	return sid.String(), nil
}
