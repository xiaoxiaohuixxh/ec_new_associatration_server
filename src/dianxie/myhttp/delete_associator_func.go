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

func httpfun_delete_associator(w http.ResponseWriter, r *http.Request) {
	//{"Delete_Associator":{"name":"8888"}}
	associator := make(map[string]interface{})

	//revbody, _ := ioutil.ReadAll(r.Body)
	//r.Body.Close()
	//fmt.Printf("%s\n", revbody)

	// json decode
	fmt.Println("客户端发过来的数据" + r.Form["data"][0])
	err := json.Unmarshal([]byte(r.Form["data"][0]), &associator)
	if err != nil {
		//panic(err)
		fmt.Println("HTTP_SERVER 识别删除会员的请求的时候出错" + err.Error())
	} else {
		associator_json := associator["Delete_Associator"].(map[string]interface{})
		var associator_message associator_pkg.Associator_s
		associator_message.Associator_Name = associator_json["Name"].(string)
		associator_message.Associator_Class = associator_json["Class"].(string)
		associator_message.Associator_Sex = associator_json["Sex"].(string)
		associator_message.Associator_Phone_Number = associator_json["Phone_Number"].(string)
		var output_zone string
		var res string
		res, err = associator_pkg.Delete_associator(associator_message)
		if err == nil {
			if res == "" {
				output_zone = "此会员输入的时候出现错误"
			} else {
				output_zone = res
			}

		} else {
			output_zone = "删除会员错误" + err.Error()
		}
		fmt.Println(associator_message.Associator_Name + "删除会员" + r.RemoteAddr)
		io.WriteString(w, output_zone)
	}

}
