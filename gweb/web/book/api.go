package book

import (
	"github.com/echopairs/skygo/gweb/model"
	"github.com/echopairs/skygo/gweb/web/commom"
	"github.com/echopairs/skygo/gweb/web/router"
	"github.com/julienschmidt/httprouter"

	"fmt"
	"net/http"
	"sync"

)

var (
	bookstore = make(map[string]*model.Book)
	mtx       sync.RWMutex
)

func init() {
	router.RegisterHttpRouteHandle("POST", "/books", "bookCreate", bookCreate)
	router.RegisterHttpRouteHandle("GET", "/books", "bookIndex", bookIndex)
	router.RegisterHttpRouteHandle("GET", "/books/:isbn", "bookShow", bookShow)
	router.RegisterHttpRouteHandle("DELETE", "/books/:isbn", "bookDelete", bookDelete)
	router.RegisterHttpRouteHandle("POST", "/books/update", "bookUpdate", bookUpdate)
}

// Handler for the books Create action
// POST /books
func bookCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	book := &model.Book{}
	err := commom.PopulateModelFromHandler(r, book)
	if err != nil {
		fmt.Printf("bookCreate PopulateModelFromHandler failed %s", err.Error())
		commom.WriteError(w, commom.ERR_INVALID_REQUEST_BODY, http.StatusBadRequest)
		return
	}
	mtx.Lock()
	defer mtx.Lock()
	if _, ok := bookstore[book.ISBN]; ok {
		fmt.Printf("bookCreate failed, book %s already exist", book.ISBN)
		commom.WriteError(w, commom.ERR_BOOK_ALREADY_EXIST, http.StatusBadRequest)
		return
	}
}

// Handler for the books index action
// GET /books
func bookIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mtx.Lock()
	defer mtx.Unlock()
	books := make([]*model.Book, len(bookstore))
	i := 0
	for _, v := range bookstore {
		books[i] = v
		i++
	}
	res := &commom.ResBody{
		Err:commom.OK,
		Data:books,
	}
	commom.WriteJson(w, res, http.StatusOK)
}

// Handler for the books Show action
// Get /books/:isbn
func bookShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isbn := params.ByName("isbn")
	mtx.Lock()
	defer mtx.Unlock()
	ok, book := isExist(isbn, w)
	if !ok {
		return
	}
	res := &commom.ResBody{
		Err:0,
		Data:book,
	}
	commom.WriteJson(w, res, http.StatusOK)
}

// Handler for delete book
// DELETE /books/:isbn
func bookDelete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isbn := params.ByName("isbn")
	mtx.Lock()
	defer mtx.Unlock()
	ok, _ := isExist(isbn, w)
	if !ok {
		return
	}
	delete(bookstore, isbn)
	commom.WriteOk(w)
}

// Handler for update book
// POST /books/update
func bookUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	book := &model.Book{}
	err := commom.PopulateModelFromHandler(r, book)
	if err != nil {
		commom.WriteError(w, commom.ERR_INVALID_REQUEST_BODY, http.StatusBadRequest)
	}
	mtx.Lock()
	defer mtx.Unlock()
	bookstore[book.ISBN] = book
	res := &commom.ResBody{
		Err:0,
		Data: bookstore[book.ISBN],
	}
	commom.WriteJson(w, res, http.StatusOK)
}

func isExist(isbn string, w http.ResponseWriter) (ok bool, book *model.Book) {
	book, ok = bookstore[isbn]
	if !ok {
		fmt.Printf("book %s not exist", isbn)
		commom.WriteError(w, commom.ERR_BOOK_NOT_EXIST, http.StatusNotFound)
		return
	}
	return
}
