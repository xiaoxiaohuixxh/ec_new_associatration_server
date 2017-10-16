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

func httpfun_change_associator_receipt_status(w http.ResponseWriter, r *http.Request) {
	//{"Change_Associator_Receipt_Status":{"name":"8888"}}
	associator := make(map[string]interface{})

	//revbody, _ := ioutil.ReadAll(r.Body)
	//r.Body.Close()
	//fmt.Printf("%s\n", revbody)

	// json decode
	fmt.Println("客户端发过来的数据" + r.Form["data"][0])
	err := json.Unmarshal([]byte(r.Form["data"][0]), &associator)
	var output_zone string
	if err != nil {
		//panic(err)
		fmt.Println("HTTP_SERVER 识别修改会员的收据状态请求的时候出错" + err.Error())
		output_zone = "识别修改会员的收据状态请求错误" + err.Error()
	} else {
		associator_json := associator["Change_Associator_Receipt_Status"].(map[string]interface{})
		var associator_message associator_pkg.Associator_s
		fmt.Println(associator_json["Name"].(string))
		associator_message.Associator_Name = associator_json["Name"].(string)
		associator_message.Associator_Class = associator_json["Class"].(string)
		associator_message.Associator_Sex = associator_json["Sex"].(string)
		associator_message.Associator_Phone_Number = associator_json["Phone_Number"].(string)

		associator_message.Receipt_Print_Status = associator_json["Receipt_Print_Status"].(string)

		timestamp := time.Now().Unix()
		timestamp_string := fmt.Sprintf("%d", timestamp)
		tm := time.Unix(timestamp, 0)
		date_string := tm.Format("2006-01-02 03:04:05 PM")
		associator_message.Receipt_Print_Server_Time = date_string
		associator_message.Receipt_Print_Server_MechineName = timestamp_string

		associator_message.Receipt_Print_Client_Time = associator_json["Print_Time"].(string)
		associator_message.Receipt_Print_Client_MechineName = associator_json["Print_MechineTime"].(string)

		associator_message.Receipt_Print_Ip = r.RemoteAddr
		associator_message.Receipt_Print_Mac = associator_json["Print_Mac"].(string)

		associator_message.Receipt_Print_Manager = associator_json["Print_Manager"].(string)

		err = associator_pkg.Change_associator_receipt_status(associator_message)
		if err == nil {
			output_zone = "修改会员的收据状态成功"
		} else {
			output_zone = "修改会员的收据状态错误：" + err.Error()
		}
		fmt.Println(associator_message.Associator_Name + "修改会员的收据状态" + r.RemoteAddr)
		io.WriteString(w, output_zone)
	}

}
