package myhttp

//20170730 注册会员http的函数已编写但是还没有测试
import (
	"dianxie/associator_pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//"time"
)

func httpfun_cancel_associator_number(w http.ResponseWriter, r *http.Request) {
	//{"Cancel_Associator_Number":{"name":"8888"}}
	associator := make(map[string]interface{})

	//revbody, _ := ioutil.ReadAll(r.Body)
	//r.Body.Close()
	//fmt.Printf("%s\n", revbody)

	// json decode
	fmt.Println("客户端发过来的数据" + r.Form["data"][0])
	err := json.Unmarshal([]byte(r.Form["data"][0]), &associator)
	if err != nil {
		//panic(err)
		fmt.Println("HTTP_SERVER 识别取消会员已申请的会员号的请求的时候出错" + err.Error())
	} else {
		associator_json := associator["Cancel_Associator_Number"].(map[string]interface{})
		var associator_message associator_pkg.Associator_s
		associator_message.Associator_Name = associator_json["Name"].(string)
		associator_message.Associator_Class = associator_json["Class"].(string)
		associator_message.Associator_Number = associator_json["Number"].(string)
		associator_message.Associator_Sex = associator_json["Sex"].(string)
		associator_message.Associator_Phone_Number = associator_json["Phone_Number"].(string)

		//var associator_number string
		var output_zone string
		_, err = associator_pkg.Cancel_associator_number(associator_message)
		if err == nil {
			output_zone = "取消会员已申请的会员号成功"
		} else {
			output_zone = "取消会员已申请的会员号错误" + err.Error()
		}
		fmt.Println(r.RemoteAddr + output_zone)
		io.WriteString(w, output_zone)
	}

}
