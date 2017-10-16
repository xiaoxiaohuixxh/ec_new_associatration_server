package myhttp

import (
	"fmt"
	"net/http"
)

var Mux map[string]func(http.ResponseWriter, *http.Request)

func Reg_url() {
	Mux = make(map[string]func(http.ResponseWriter, *http.Request))
	Mux["/hello/"] = Hello
	Mux["/bye/"] = Bye
	Mux["/page/"] = Page
	fmt.Println("reg url ok")
}
