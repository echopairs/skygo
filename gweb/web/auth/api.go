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
	router.RegisterHttpRouteHandle("GET", "/users/:name", "getUser", CheckPrivilegesWithRouterHandle(getUser, "getUser"))
	router.RegisterHttpHandleFunc("POST", "/users", "createUser", CheckPrivilegesWithHttpHandle(createUser, "createUser"))
	router.RegisterHttpHandleFunc("GET", "/roles", "getAllRoles", CheckPrivilegesWithHttpHandle(getAllRoles, "getAllRoles"))
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
	var access []model.Access

	//sqlStr := "select access_id from role_access where role_id = (select role_id from user_role where user_id = ?)"
	sqlStr := `select access.id, access.name from user_role as ur, role_access as ra, access where 
				ur.user_id = ?
				AND ur.role_id = ra.role_id
				AND ra.access_id = access.id
				`
	err = authDb.Select(&access, sqlStr, user.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			resBody.Err = common.ERR_EXEC_QUERY_SQL_ERROR
			resBody.Msg = "get access  failed"
			common.WriteJson(w, resBody, http.StatusInternalServerError)
			log.Println(err)
			return
		}
		log.Printf("get access failed %s\n", err.Error())
	}

	for _, value := range access {
		user.Roles = append(user.Roles, value.Name)
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
	resBody.Data = access
	common.WriteJson(w, resBody, http.StatusOK)
}

func logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("logout")
	sess.SessionDestroy(w, r)
	common.WriteOk(w)
}

// GET /users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []*model.User
	err := authDb.Select(&users, "select * from user")
	if err != nil {
		log.Printf("query sql err %s", err.Error())
		common.WriteError(w, common.ERR_EXEC_QUERY_SQL_ERROR, http.StatusInternalServerError)
		return
	}
	resBody := &common.ResBody{
		Err:  common.OK,
		Data: users,
	}
	common.WriteJson(w, resBody, http.StatusOK)
	return
}

// GET /users/:name
func getUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var user model.User
	resBody := &common.ResBody{}
	userName := params.ByName("name")
	log.Printf("userName is %s", userName)
	err := authDb.Get(&user, "select * from user where name = ?", userName)
	if err != nil {
		log.Printf("get user query err %s ", err.Error())
		if err != sql.ErrNoRows {
			resBody.Err = common.ERR_EXEC_QUERY_SQL_ERROR
			resBody.Msg = common.GetError(common.ERR_EXEC_QUERY_SQL_ERROR) + err.Error()
			common.WriteJson(w, resBody, http.StatusInternalServerError)
			return
		}
		resBody.Err = common.ERR_QUERY_USER_NOT_EXIST
		resBody.Msg = common.GetError(common.ERR_QUERY_USER_NOT_EXIST) + ":" + userName
		common.WriteJson(w, resBody, http.StatusNotFound)
		return
	}

	// 2. get role by uid
	var access []model.Access
	sqlStr := `select access.id, access.name from user_role as ur, role_access as ra, access where
				ur.user_id = ?
				AND ur.role_id = ra.role_id
				AND ra.access_id = access.id`
	err = authDb.Select(&access, sqlStr, user.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			resBody.Err = common.ERR_EXEC_QUERY_SQL_ERROR
			resBody.Msg = "get access  failed"
			common.WriteJson(w, resBody, http.StatusInternalServerError)
			log.Println(err)
			return
		}
		log.Printf("get access failed %s\n", err.Error())
	}

	for _, value := range access {
		user.Roles = append(user.Roles, value.Name)
	}
	resBody.Err = common.OK
	resBody.Data = user
	common.WriteJson(w, resBody, http.StatusOK)
}

// POST /users
func createUser(w http.ResponseWriter, r *http.Request) {
	body := common.ParseForm(r)
	resBody := &common.ResBody{}
	if body.GetErr() != nil {
		resBody.Err = common.ERR_INVALID_REQUEST_PARAMS
		resBody.Msg = common.GetError(common.ERR_INVALID_REQUEST_PARAMS)
		common.WriteJson(w, resBody, http.StatusBadRequest)
		return
	}
	username, _ := body.GetStringVar("username")
	password, _ := body.GetStringVar("password")
	rolename, _ := body.GetStringVar("rolename")
	log.Printf(username, password, rolename)
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

// GET /roles
func getAllRoles(w http.ResponseWriter, r *http.Request) {
	var roles []*model.Role
	resBody := &common.ResBody{
		Err: common.OK,
	}
	err := authDb.Select(&roles, "select *from role")
	if err != nil {
		if err == sql.ErrNoRows {
			common.WriteOk(w)
			return
		}
		resBody.Err = common.ERR_EXEC_QUERY_SQL_ERROR
		resBody.Msg = err.Error()
		common.WriteJson(w, resBody, http.StatusInternalServerError)
	}
	resBody.Data = roles
	common.WriteJson(w, roles, http.StatusOK)
}