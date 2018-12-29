package auth

import (
	"github.com/echopairs/skygo/gweb/web/router"
	"net/http"
)

func init() {
	router.RegisterHttpHandleFunc("POST", "/login", "login", login)
	router.RegisterHttpHandleFunc("POST", "/logout", "logout", logout)
}

func login(w http.ResponseWriter, r *http.Request) {
	// todo
}

func logout(w http.ResponseWriter, r *http.Request) {
	// todo
}
