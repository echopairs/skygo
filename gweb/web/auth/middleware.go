package auth

import (
	"net/http"
)

func CheckPrivileges(h http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. check privileges here

		// 2. after privileges check
		h(w, r)
	}
}
