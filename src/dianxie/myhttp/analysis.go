package myhttp

import (
	"fmt"
	"html"
	"io"
	"net/http"

	"strings"
)

type MyHandle struct{}

func (*MyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := Mux[r.URL.String()]; ok {
		h(w, r)
	}
	//io.WriteString(w, "URL"+r.URL.String())
	//fmt.Println("HTTP SERVER in")
	r.ParseForm() //解析参数,默认不不解析的
	//fmt.Println(r.Method)
	//io.WriteString(w, "HTTP SERVER 模块")
	if r.Method == "GET" {
		if strings.Contains(r.URL.String(), "/PC_APP_API/Register_Associator") {
			fmt.Println("Http_server Register_Associator")
			httpfun_register_associator(w, r)
		} else {
			fmt.Println("Http_server SERVER CALL BACK REV")
		}

	} else if r.Method == "POST" {
		fmt.Println("HTTP SERVER ASK API REV")
		if strings.Contains(r.URL.String(), "/PC_APP_API/Register_Associator") {
			//注册会员
			fmt.Println("Http_server Register_Associator")
			httpfun_register_associator(w, r)
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Get_Associator_List") {
			//获取会员列表
			fmt.Println("Http_server Get_Associator_List")
			httpfun_get_associator_list(w, r)
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Get_Associator_Number") {
			//申请会员的会员号
			fmt.Println("Http_server Get_Associator_Number")
			httpfun_get_associator_number(w, r)
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Change_Associator_Receipt_Status") {
			//修改会员的收据打印状态
			fmt.Println("Http_server Change_Associator_Receipt_Status")
			httpfun_change_associator_receipt_status(w, r)
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Change_Associator_Card_Id") {
			//修改会员的会员卡号
			fmt.Println("Http_server Change_Associator_Card_Id")
			httpfun_change_associator_card_id(w, r)
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Cancel_Associator_Number") {
			//取消会员已申请的会员号
			fmt.Println("Http_server Cancel_Associator_Number")
			httpfun_cancel_associator_number(w, r)
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Check_Associator_Exist_Card_Number") {
			//查询此会员是否存在会员卡
			fmt.Println("Http_server Check_Associator_Exist_Card_Number")
			httpfun_check_associator_exist_card_number(w, r)
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Login_Administrator") {
			//管理员登录
			fmt.Println("Http_server Login_Administrator")
			httpfun_login_administrator(w, r)
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Check_Associator_Receipt_Status") {
			//查询此会员当前收据打印状态
			fmt.Println("Http_server Check_Associator_Receipt_Status")
			httpfun_check_associator_receipt_status(w, r)
			//
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Check_Associator_Receipt_Adminstartor") {
			//查询此会员当前收据打印的管理员
			fmt.Println("Http_server Check_Associator_Receipt_Adminstartor")
			httpfun_check_associator_receipt_adminstrator(w, r)
			//
		} else if strings.Contains(r.URL.String(), "/PC_APP_API/Delete_Associator") {
			//查询此会员当前收据打印的管理员
			fmt.Println("Http_server Delete_Associator")
			httpfun_delete_associator(w, r)
			//
		} else {
			fmt.Println("Http_server SERVER CALL BACK REV")
		}
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello 模块")
}

func Bye(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "bye 模块")
}

func Page(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数,默认不不解析的
	if r.Method == "GET" {
		fmt.Sprintln(w, "   enen hi", r.Form["hi"])
	}
	io.WriteString(w, "hello 模块"+html.EscapeString(r.URL.Path[1:]))
}
