package auth

import (
	"log"
	"net/http"

	"github.com/echopairs/skygo/gweb/web/common"
	"github.com/julienschmidt/httprouter"
)

func CheckPrivilegesWithHttpHandle(h http.HandlerFunc, funcname string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. get user first
		user := GetUser(r)
		resBody := common.ResBody{}
		if user == nil {
			resBody.Err = common.ERR_USER_NOT_LOGIN
			resBody.Msg = common.GetError(common.ERR_USER_NOT_LOGIN)
			common.WriteJson(w, resBody, http.StatusBadRequest)
			return
		}

		ok, errCode := checkPrivileges(user.ID, funcname)
		if !ok {
			common.WriteError(w, errCode, http.StatusForbidden)
			return
		}
		// 3. after privileges check
		h(w, r)
	}
}

func CheckPrivilegesWithRouterHandle(h httprouter.Handle, funcname string) httprouter.Handle {
	return nil
}

func checkPrivileges(id int, funcname string) (bool, int) {

	sql := `select count(ur.user_id) from user_role ur, role_access ra, access a where
			ur.user_id = ?
			AND ra.role_id = ur.role_id
			AND ra.access_id = a.id
			AND a.name = ?`
	count := 0
	err := authDb.Get(&count, sql, id, funcname)
	if err != nil {
		log.Printf("select count failed %s ", err.Error())
		return false, common.ERR_EXEC_QUERY_SQL_ERROR
	}

	if count == 0 {
		log.Printf("err auth unauthorized\n")
		return false, common.ERR_AUTH_UNAUTHORIZED
	}

	return true, common.OK
}
