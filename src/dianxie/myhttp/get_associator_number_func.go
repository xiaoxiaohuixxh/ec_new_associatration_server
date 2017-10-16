package myhttp

//20170730 注册会员http的函数已编写但是还没有测试
import (
	"dianxie/associator_pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func httpfun_get_associator_number(w http.ResponseWriter, r *http.Request) {
	//{"Get_Associator_Number":{"name":"8888"}}
	associator := make(map[string]interface{})
	//revbody := ""
	//revbody, _ := ioutil.ReadAll(r.Body)
	//r.Body.Close()
	//fmt.Printf("%s\n", revbody)

	// json decode
	var output_zone string
	fmt.Println(r.Form["data"][0])
	err := json.Unmarshal([]byte(r.Form["data"][0]), &associator)
	if err != nil {
		//panic(err)
		fmt.Println("HTTP_SERVER 识别申请会员的会员号的请求的时候出错" + err.Error())
		output_zone = "申请会员的会员号错误"
	} else {
		associator_json := associator["Get_Associator_Number"].(map[string]interface{})
		var associator_message associator_pkg.Associator_s
		associator_message.Associator_Name = associator_json["Name"].(string)
		associator_message.Associator_Class = associator_json["Class"].(string)
		associator_message.Associator_Sex = associator_json["Sex"].(string)
		associator_message.Associator_Phone_Number = associator_json["Phone_Number"].(string)

		associator_message.Associator_Number, err = associator_pkg.Get_associator_number(associator_message)

		if err != nil {
			output_zone = "申请会员的会员号错误"
		} else {
			var associator_json_text []byte
			associator_list_json := make(map[string]interface{})
			associator_list_json["Get_Associator_Number"] = associator_message
			associator_json_text, _ = json.Marshal(associator_list_json)
			output_zone = string(associator_json_text)
		}

		io.WriteString(w, output_zone)
	}

}
