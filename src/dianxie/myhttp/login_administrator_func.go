package myhttp

import (
	"dianxie/associator_pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//"time"
)

func httpfun_login_administrator(w http.ResponseWriter, r *http.Request) {
	//{"Login_Administrator":{"name":"8888"}}
	associator := make(map[string]interface{})

	//revbody, _ := ioutil.ReadAll(r.Body)
	//r.Body.Close()
	//fmt.Printf("%s\n", revbody)

	// json decode
	fmt.Println("客户端发过来的数据" + r.Form["data"][0])
	err := json.Unmarshal([]byte(r.Form["data"][0]), &associator)
	if err != nil {
		//panic(err)
		fmt.Println("HTTP_SERVER 管理员登录的请求的时候出错" + err.Error())
	} else {
		associator_json := associator["Login_Administrator"].(map[string]interface{})
		var administrator_name string
		var administrator_passwd string
		administrator_name = associator_json["Name"].(string)
		administrator_passwd = associator_json["Passwd"].(string)
		var output_zone string
		var res bool
		res, err = associator_pkg.Login_administrator(administrator_name, administrator_passwd)
		if err == nil {
			if res {
				output_zone = "登录成功"
			} else {
				output_zone = "登录失败"
			}

		} else {
			output_zone = "登录失败" + err.Error()
		}
		fmt.Println(administrator_name + "登录" + r.RemoteAddr)
		io.WriteString(w, output_zone)
	}

}
