package main

import (
	"github.com/ddliu/go-httpclient"
	"io/ioutil"

	"fmt"
	"net/http"
)

var loginString string = `{
	"username": "root",
	"password": "admin"
}`

func main() {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "my awsome httpclient",
		"Accept-Language":        "en-us",
	})

	res, err := httpclient.PostJson("http://127.0.0.1:9090/login", loginString)
	if err != nil {
		fmt.Errorf("error %s", err.Error())
	}
	//var body interface{}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("res is %s\n", string(data))

	cooks := res.Cookies()
	sname := cooks[0].Name
	svalue := cooks[0].Value

	fmt.Printf("for test getAllUsers\n")
	res, err = httpclient.WithCookie(
		&http.Cookie{
			Name: sname,
			Value: svalue,
		}).Get("http://127.0.0.1:9090/users")

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("/users response is %s\n", string(data))

	////fmt.Println(res)
	//res, err = httpclient.WithCookie(
	//	&http.Cookie{
	//		Name:  sname,
	//		Value: svalue,
	//	}).PostJson("http://127.0.0.1:9090/logout", loginString)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
}
