package commom

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// deal with http response

type ResBody struct {
	Err  int         `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func WriteJson(w http.ResponseWriter, v interface{}, httpCode int) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("Invalid json %v", v)
		msg := fmt.Sprintf("{\"err\":500,\"msg\":\"%s\"}", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(data)
}

func WriteOk(w http.ResponseWriter) {
	WriteJson(w, &ResBody{Err: 0, Msg: "success"}, http.StatusOK)
}

func WriteError(w http.ResponseWriter, errCode int, httpCode int) {
	v := &ResBody{
		errCode,
		GetError(errCode),
		httpCode,
	}
	WriteJson(w, v, httpCode)
}
