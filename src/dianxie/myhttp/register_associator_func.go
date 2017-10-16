package myhttp

//20170730 注册会员http的函数已编写但是还没有测试
import (
	"dianxie/associator_pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func httpfun_register_associator(w http.ResponseWriter, r *http.Request) {
	//{"Register_Associator":{"name":"8888"}}
	associator := make(map[string]interface{})
	revbody := ""
	//revbody, _ := ioutil.ReadAll(r.Body)
	//r.Body.Close()
	fmt.Printf("%s\n", revbody)

	// json decode
	fmt.Println(r.Form["data"][0])
	err := json.Unmarshal([]byte(r.Form["data"][0]), &associator)
	var output_zone string
	if err != nil {
		//panic(err)
		fmt.Println("HTTP_SERVER 识别会员注册的请求的时候出错" + err.Error())
		output_zone = "识别会员注册的请求的时候错误"
	} else {
		associator_json := associator["Register_Associator"].(map[string]interface{})
		var associator_message associator_pkg.Associator_s
		associator_message.Associator_Name = associator_json["Name"].(string)
		associator_message.Associator_Class = associator_json["Class"].(string)
		associator_message.Associator_Sex = associator_json["Sex"].(string)
		associator_message.Associator_Phone_Number = associator_json["Phone_Number"].(string)
		associator_message.Associator_QQ_Number = associator_json["QQ_Number"].(string)
		associator_message.Associator_Wechat_Number = associator_json["Wechat_Number"].(string)

		timestamp := time.Now().Unix()
		timestamp_string := fmt.Sprintf("%d", timestamp)
		tm := time.Unix(timestamp, 0)
		date_string := tm.Format("2006-01-02 03:04:05 PM")
		associator_message.Associator_Register_Server_Time = date_string
		associator_message.Associator_Register_Server_MechineTime = timestamp_string

		associator_message.Associator_Register_Client_Time = associator_json["Register_Time"].(string)
		associator_message.Associator_Register_Client_MechineTime = associator_json["Register_MechineTime"].(string)

		associator_message.Associator_Register_Ip = r.RemoteAddr
		associator_message.Associator_Register_Mac = associator_json["Register_Mac"].(string)
		err = associator_pkg.Register_associator(associator_message)

		if err == nil {
			output_zone = "注册会员成功"
		} else {
			output_zone = "注册会员错误" + err.Error()
		}

	}
	io.WriteString(w, output_zone)
}
