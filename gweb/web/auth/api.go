package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/echopairs/skygo/gweb/model"
	"github.com/echopairs/skygo/gweb/web/common"
	"github.com/echopairs/skygo/gweb/web/router"
	"github.com/julienschmidt/httprouter"
)

func init() {
	router.RegisterHttpHandleFunc("POST", "/login", "login", login)
	router.RegisterHttpHandleFunc("POST", "/logout", "logout", logout)
	router.RegisterHttpHandleFunc("GET", "/users", "getAllUsers", CheckPrivilegesWithHttpHandle(getAllUsers, "getAllUsers"))
}

func login(w http.ResponseWriter, r *http.Request) {
	// 1. verify
	body := common.ParseForm(r)
	resBody := &common.ResBody{}
	if body.GetErr() != nil {
		resBody.Err = common.ERR_INVALID_REQUEST_BODY
		resBody.Msg = common.GetError(common.ERR_INVALID_REQUEST_BODY)
		common.WriteJson(w, resBody, http.StatusBadRequest)
		return
	}
	username, err := body.GetStringVar("username")
	if err != nil {
		resBody.Err = common.ERR_INVALID_REQUEST_PARAMS
		resBody.Msg = common.GetError(common.ERR_INVALID_REQUEST_PARAMS)
		common.WriteJson(w, resBody, http.StatusBadRequest)
		return
	}
	password, err := body.GetStringVar("password")

	if err != nil {
		resBody.Err = common.ERR_INVALID_REQUEST_PARAMS
		resBody.Msg = common.GetError(common.ERR_INVALID_REQUEST_PARAMS)
		common.WriteJson(w, resBody, http.StatusBadRequest)
	}

	user := model.User{}
	err = authDb.Get(&user, "select id, name, password, salt from user where name = ?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			resBody.Err = common.ERR_INVALID_AUTH_USER_NOT_EXIST
			resBody.Msg = common.GetError(common.ERR_INVALID_AUTH_USER_NOT_EXIST)
			common.WriteJson(w, resBody, http.StatusBadRequest)
			return
		}
		resBody.Err = common.ERR_EXEC_QUERY_SQL_ERROR
		resBody.Msg = common.GetError(common.ERR_EXEC_QUERY_SQL_ERROR) + err.Error()
		common.WriteJson(w, resBody, http.StatusInternalServerError)
		log.Printf("error %s", err.Error())
		return
	}
	if !user.VerifyPassword(password) {
		resBody.Err = common.ERR_INVALID_AUTH_PASSWORD
		resBody.Msg = common.GetError(common.ERR_INVALID_AUTH_PASSWORD)
		return
	}

	// 2. get roles
	var accessId []int
	sqlStr := "select access_id from role_access where role_id = (select role_id from user_role where user_id = ?)"
	err = authDb.Select(&accessId, sqlStr, user.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			resBody.Err = common.ERR_EXEC_QUERY_SQL_ERROR
			resBody.Msg = "get access ids failed"
			common.WriteJson(w, resBody, http.StatusInternalServerError)
			log.Println(err)
			return
		}
		log.Printf("get access ids failed %s\n", err.Error())
	}

	var roles []string
	tx, err := authDb.Beginx()
	if err != nil {
		log.Println(err)
		return
	}
	st, err := tx.Preparex("select name from access where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	var role string
	for _, id := range accessId {
		if err = st.Get(&role, id); err != nil {
			// todo
			log.Println(err)
			return
		}
		roles = append(roles, role)
	}
	tx.Commit()

	for _, role = range roles {
		user.Roles = append(user.Roles, role)
	}

	// 3. write to session
	sess, err := sess.SessionCreate(w, r)
	if err != nil {
		resBody.Err = common.ERR_CREATE_SESSION_ERROR
		resBody.Msg = common.GetError(common.ERR_CREATE_SESSION_ERROR) + err.Error()
		common.WriteJson(w, resBody, http.StatusInternalServerError)
		return
	}

	if err = sess.Set(userKey, &user); err != nil {
		resBody.Err = common.ERR_SET_USER_TO_SESSION_ERROR
		resBody.Msg = common.GetError(common.ERR_SET_USER_TO_SESSION_ERROR) + err.Error()
		common.WriteJson(w, resBody, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	resBody.Data = user.Roles
	common.WriteJson(w, resBody, http.StatusOK)
}

func logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("logout")
	sess.SessionDestroy(w, r)
	common.WriteOk(w)
}

// GET /users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	var users [] *model.User
	err := authDb.Select(&users, "select * from user")
	if err != nil {
		log.Printf("query sql err %s", err.Error())
		common.WriteError(w, common.ERR_EXEC_QUERY_SQL_ERROR, http.StatusInternalServerError)
		return
	}
	resBody := &common.ResBody{
		Err: common.OK,
		Data: users,
	}
	common.WriteJson(w, resBody, http.StatusOK)
	return
}

// GET /users/:name
func getUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}

// POST /users
func createUser(w http.ResponseWriter, r *http.Request) {

}

// DELETE /users/:name
func deleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}

// PUT /users/:name/password
func updateUserPassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}

// PUT /users/:name/roles
func updateUserRoles(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}
