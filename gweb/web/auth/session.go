package auth

import (
	"github.com/echopairs/skygo/session"
)

type HttpConfig = session.ManagerConfig

type SessionStorage struct {
	*session.Manager
}

func NewSessionStorage(cfg *HttpConfig) (*SessionStorage, error) {
	manager, err := session.NewManager(cfg)
	return &SessionStorage{
		manager,
	}, err
}
