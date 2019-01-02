package common

const (
	OK = 0

	// 1-99: common errors
	ERR_INVALID_REQUEST_PARAMS = 1
	ERR_INVALID_REQUEST_BODY   = 2
	ERR_INVALID_AUTH_USERNAME = 3
	ERR_INVALID_AUTH_PASSWORD = 4
	ERR_INVALID_AUTH_USER_NOT_EXIST = 5
	ERR_EXEC_QUERY_SQL_ERROR = 6
	ERR_CREATE_SESSION_ERROR = 7
	ERR_SET_USER_TO_SESSION_ERROR = 8


	// 100-199 Books
	ERR_BOOK_CREATE_FAILED = 101
	ERR_BOOK_QUERY_FAILED  = 102
	ERR_BOOK_NOT_EXIST     = 103
	ERR_BOOK_ALREADY_EXIST = 104
)

var errMap = map[int]string{
	ERR_INVALID_REQUEST_PARAMS: "Invalid request param",
	ERR_INVALID_REQUEST_BODY:   "Invalid request body",
	ERR_INVALID_AUTH_USERNAME:	"Invalid request username",
	ERR_INVALID_AUTH_USER_NOT_EXIST: "Invalid user not exist",
	ERR_INVALID_AUTH_PASSWORD: "Invalid user password",
	ERR_EXEC_QUERY_SQL_ERROR: "Exec sql error",
	ERR_CREATE_SESSION_ERROR: "Create session failed ",
	ERR_SET_USER_TO_SESSION_ERROR: "Set user to session failed ",

	ERR_BOOK_CREATE_FAILED: "Book create failed",
	ERR_BOOK_QUERY_FAILED:  "Book query failed",
	ERR_BOOK_NOT_EXIST:     "Book not exist",
	ERR_BOOK_ALREADY_EXIST:  "Book already exist",
}

func GetError(code int) string {
	s, ok := errMap[code]
	if ok {
		return s
	}
	return "Unknown error"
}
