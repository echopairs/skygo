package commom

const (
	OK = 0

	// 1-99: common errors
	ERR_INVALID_REQUEST_PARAMS = 1
	ERR_INVALID_REQUEST_BODY   = 2

	// 100-199 Books
	ERR_BOOK_CREATE_FAILED = 101
	ERR_BOOK_QUERY_FAILED  = 102
	ERR_BOOK_NOT_EXIST     = 103
	ERR_BOOK_ALREADY_EXIST = 104
)

var errMap = map[int]string{
	ERR_INVALID_REQUEST_PARAMS: "Invalid Request param",
	ERR_INVALID_REQUEST_BODY:   "Invalid request Body",

	ERR_BOOK_CREATE_FAILED: "Book create Failed",
	ERR_BOOK_QUERY_FAILED:  "Book query Failed",
	ERR_BOOK_NOT_EXIST:     "Book Not Exist",
	ERR_BOOK_ALREADY_EXIST:  "Book Already Exist",
}

func GetError(code int) string {
	s, ok := errMap[code]
	if ok {
		return s
	}
	return "Unknown error"
}
