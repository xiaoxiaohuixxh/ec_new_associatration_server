package main

import (
	"dianxie/appconf"
	"dianxie/associator_pkg"
	"dianxie/myhttp"
	"dianxie/mylog"
	"dianxie/sent_udp_heartbeak"
	//"encoding/json"
	"fmt"
	"sync"
)

var W sync.WaitGroup

func main() {
	mylog.StartLog()
	mylog.ErrorLog("associator_system_start\r\n")
	appconf.Read_conf()
	fmt.Println("-------------------------------------")
	fmt.Println("SQLLITE3地址:", appconf.Tf.SqlDb)
	fmt.Println("SQLLITE3会员数据表名:", appconf.Tf.SqlAssociatorTable)
	fmt.Println("SQLLITE3会员会员号数据表名:", appconf.Tf.SqlAssociatorNumberTable)
	fmt.Println("-------------------------------------")
	fmt.Println("HTTP服务端地址:", appconf.Tf.HttpAddr)
	fmt.Println("HTTP服务端使能HTTPS:", appconf.Tf.HttpsEnable)
	fmt.Println("HTTP服务端HttpsCertFile路径:", appconf.Tf.HttpsCertFile)
	fmt.Println("HTTP服务端HttpsKeyFile路径:", appconf.Tf.HttpsKeyFile)
	fmt.Println("HTTP服务端读超时时间:", appconf.Tf.HttpReadTimeout)
	fmt.Println("HTTP服务端写超时时间:", appconf.Tf.HttpWriteTimeout)
	fmt.Println("HTTP服务端最大头大小:", appconf.Tf.HttpMaxHeaderBytes)
	fmt.Println("-------------------------------------")
	//数据库初始化
	associator_pkg.Init_sqllite_databse()
	W.Add(1)

	//测试修改收据
	/*
		var associator associator_pkg.Associator_s
		associator.Associator_Name = "123"
		associator.Associator_Class = "123"
		associator.Associator_Sex = "123"
		associator.Associator_Phone_Number = "123"
		associator_pkg.Register_associator(associator)

		associator.Receipt_Print_Status = "proceed"
		associator_pkg.Change_associator_receipt_status(associator)

		associator_pkg.Get_associator_number(associator)

		associator_pkg.Cancel_associator_number(associator)

		associator.Receipt_Print_Status = "no_proceed"
		associator_pkg.Change_associator_receipt_status(associator)

		associator.Associator_Card_Number = "8888888888"
		associator_pkg.Change_associator_card_id(associator)

		associator_list_table, err := associator_pkg.Get_associator_list_by_receipt_status("no_proceed")
		if err != nil {
			fmt.Println("sqlite3数据库获取会员列表err:" + err.Error())
		} else {

			var associator_json_text []byte
			associator_list_json := make(map[string]interface{})
			associator_list_json["Associator_List"] = associator_list_table
			associator_json_text, _ = json.Marshal(associator_list_json)

			fmt.Println("sqlite3数据库获取会员列表" + string(associator_json_text))
		}
	*/
	//Http服务初始化
	go myhttp.Http_server_init()
	//Udp心跳广播服务初始化
	go sent_udp_heartbeak.Boardcast_udp_heartbeak_to_all_netcard()
	//	for {
	//	}
	W.Wait()
	associator_pkg.Db.Close()
}
