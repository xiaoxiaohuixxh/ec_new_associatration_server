package myhttp

//20170730 注册会员http的函数已编写但是还没有测试
import (
	"dianxie/associator_pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func httpfun_get_associator_list(w http.ResponseWriter, r *http.Request) {
	//{"Get_Associator_List":{"name":"8888"}}
	associator := make(map[string]interface{})
	var output_zone string
	//revbody, _ := ioutil.ReadAll(r.Body)
	//r.Body.Close()
	//fmt.Printf("%s\n", revbody)

	// json decode
	fmt.Println("客户端发过来的数据" + r.Form["data"][0])
	err := json.Unmarshal([]byte(r.Form["data"][0]), &associator)
	if err != nil {
		//panic(err)
		fmt.Println("HTTP_SERVER 识别修改会员的收据状态请求的时候出错" + err.Error())
		output_zone = "sqlite3数据库获取会员列表错误err:" + err.Error()
	} else {
		associator_json := associator["Get_Associator_List"].(map[string]interface{})
		var associator_message associator_pkg.Associator_s

		associator_message.Receipt_Print_Status = associator_json["Receipt_Print_Status"].(string)

		associator_list_table, err := associator_pkg.Get_associator_list_by_receipt_status(associator_message.Receipt_Print_Status)
		if err != nil {
			fmt.Println("sqlite3数据库获取会员列表错误err:" + err.Error())
			output_zone = "sqlite3数据库获取会员列表错误err:" + err.Error()
		} else {

			var associator_json_text []byte
			associator_list_json := make(map[string]interface{})
			associator_list_json["Associator_List"] = associator_list_table
			associator_json_text, _ = json.Marshal(associator_list_json)
			output_zone = string(associator_json_text)
			fmt.Println("sqlite3数据库获取会员列表" + string(len(associator_list_table)) + string(associator_json_text))
		}
		fmt.Println(associator_message.Associator_Name + "修改会员的收据状态" + r.RemoteAddr)

	}
	io.WriteString(w, output_zone)

}
