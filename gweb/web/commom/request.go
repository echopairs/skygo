package commom

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// deal with http request

const (
	ContentTypeJson   = 0
	ContentUrlEncoded = 1
	ContentTypeOther  = 2
)

type ReqBody struct {
	contentType int
	r           *http.Request
	err         error
	jsonMap     map[string]interface{}
	formMap     url.Values
}

func PopulateModelFromHandler(r *http.Request, model interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}

func ParseForm(r *http.Request) *ReqBody {
	req := &ReqBody{
		contentType: ContentTypeJson,
		r:           r,
	}

	// 1. check
	if r == nil || r.Body == nil {
		req.err = errors.New("body is nil")
		return req
	}

	// 2. read from body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		req.err = err
		return req
	}

	// 3. parse
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "application/json") {
		parseJson(b, req)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		parseUrlEncoded(b, req)
	} else {
		req.contentType = ContentTypeOther
		req.err = fmt.Errorf("Invalid Content-Type %s ", contentType)
	}
	return req
}

func parseJson(input []byte, req *ReqBody) {
	req.contentType = ContentTypeJson
	ob := make(map[string]interface{})
	err := json.Unmarshal(input, &ob)
	if err != nil {
		log.Fatalf("parse json error %v", err)
		req.err = err
		return
	} else {
		req.jsonMap = ob
	}
}

func parseUrlEncoded(input []byte, req *ReqBody) {
	req.contentType = ContentUrlEncoded
	form, err := url.ParseQuery(string(input))
	if err != nil {
		log.Fatalf("parse url encoded error %v", err)
		req.err = err
		return
	} else {
		req.formMap = form
	}
}

func (req *ReqBody) GetStringVar(key string) (value string, err error) {
	if req.err != nil {
		err = req.err
		return
	}
	if req.contentType == ContentTypeJson {
		v, ok := req.jsonMap[key]
		if !ok {
			err = fmt.Errorf("jsonMap there is no this %s field ", key)
			return
		}
		value, ok = v.(string)
		if !ok {
			err = fmt.Errorf("jsonmap type change to string failed")
			return
		}

	} else {
		// urlEncoded
		v, ok := req.formMap[key]
		if !ok {
			err = fmt.Errorf("formMap there is no this %s field ", key)
			return
		}
		value = v[0]
	}
	return
}

func (req *ReqBody) GetIntVar(key string) (value int, err error) {
	str, err := req.GetStringVar(key)
	if err != nil {
		return
	}
	value, err = strconv.Atoi(str)
	return
}

func (req *ReqBody) GetBoolVar(key string) (value bool, err error) {
	if req.err != nil {
		err = req.err
		return
	}
	if req.contentType == ContentTypeJson {
		v, ok := req.jsonMap[key]
		if !ok {
			err = fmt.Errorf("jsonMap there is no this  %s field", key)
			return
		}
		value, ok = v.(bool)
		if !ok {
			err = fmt.Errorf("type change %s field to bool failed", key)
			return
		}
	} else {
		v, ok := req.formMap[key]
		if !ok {
			err = fmt.Errorf("formMap there is no this %s field", key)
			return
		}
		if v[0] == "true" {
			value = true
		} else {
			value = false
		}
	}
	return
}

func (req *ReqBody) GetErr() error {
	return req.err
}
