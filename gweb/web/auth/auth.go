package auth

import (
	"github.com/echopairs/skygo/gweb/model"
	"github.com/echopairs/skygo/zsql"

	"fmt"
	"net/http"
)

type ctxKey struct {
	name string
}

var (
	sess      *SessionStorage
	sessionId = "gwebId"
	userKey   = &ctxKey{"user"}

	authDb *zsql.DB
)

func Set(ss *SessionStorage, db *zsql.DB) {
	sess = ss
	authDb = db
}

func GetUser(r *http.Request) *model.User {
	// 1. check context first
	ctx := r.Context()
	v := ctx.Value(userKey)
	if v != nil {
		return v.(*model.User)
	}

	// 2. load from session
	if sess == nil {
		fmt.Printf("please set sess first")
		return nil
	}
	sid, err := r.Cookie(sessionId)
	if err != nil {
		return nil
	}
	if item, ok := sess.SessionRead(sid.Value); ok != nil {
		// not found
		return nil
	} else {
		u, ok := item.Get(userKey).(*model.User)
		if ok {
			return u
		} else {
			return nil
		}
	}

	return nil
}
