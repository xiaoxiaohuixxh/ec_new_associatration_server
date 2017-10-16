package associator_pkg

import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Get_associator_list_by_receipt_status(receipt_Print_Status string) ([]interface{}, error) {
	//	var result_json string
	var associator_list_table []interface{}

	Sqlite_operation_lock.Lock()
	sql_text := "SELECT"
	sql_text += " Associator_Name,"
	sql_text += " Associator_Number,"
	sql_text += " Associator_Class,"
	sql_text += " Associator_Sex,"
	sql_text += " Associator_Phone_Number,"
	sql_text += " Associator_Card_Number,"
	sql_text += " Associator_QQ_Number,"
	sql_text += " Associator_Wechat_Number,"
	sql_text += " Associator_Register_Server_Time,"
	sql_text += " Associator_Register_Server_MechineTime,"
	sql_text += " Associator_Register_Client_Time,"
	sql_text += " Associator_Register_Client_MechineTime,"
	sql_text += " Receipt_Print_Manager"
	sql_text += " FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Receipt_Print_Status=?"

	stmt, err := Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库获取会员列表sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库获取会员列表sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		Sqlite_operation_lock.Unlock()
		return associator_list_table, err

	}
	var query *sql.Rows
	query, err = stmt.Query(receipt_Print_Status)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库获取会员列表查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库获取会员列表查询错误(err:" + err.Error() + ")")
	} else {

		for query.Next() {
			fmt.Println("sqlite3数据库获取会员列表会员号中")
			var associator_Name string
			var associator_Number string
			var associator_Class string
			var associator_Sex string
			var associator_Phone_Number string
			var associator_Card_Number string
			var associator_QQ_Number string
			var associator_Wechat_Number string
			var associator_Register_Server_Time string
			var associator_Register_Server_MechineTime string
			var associator_Register_Client_Time string
			var associator_Register_Client_MechineTime string
			var associator_Receipt_Print_Manager string
			query.Scan(
				&associator_Name,
				&associator_Number,
				&associator_Class,
				&associator_Sex,
				&associator_Phone_Number,
				&associator_Card_Number,
				&associator_QQ_Number,
				&associator_Wechat_Number,
				&associator_Register_Server_Time,
				&associator_Register_Server_MechineTime,
				&associator_Register_Client_Time,
				&associator_Register_Client_MechineTime,
				&associator_Receipt_Print_Manager)
			if associator_Name == "" {
				err = errors.New("会员数据库出现严重错误！！")

			} else {
				associator_json := make(map[string]interface{})
				associator_json["Name"] = associator_Name
				associator_json["Number"] = associator_Number
				associator_json["Class"] = associator_Class
				associator_json["Sex"] = associator_Sex
				associator_json["Phone"] = associator_Phone_Number
				associator_json["Card"] = associator_Card_Number
				associator_json["Qq"] = associator_QQ_Number
				associator_json["Wechat"] = associator_Wechat_Number
				associator_json["Register_Server_Time"] = associator_Register_Server_Time
				associator_json["Register_Server_MechineTime"] = associator_Register_Server_MechineTime
				associator_json["Register_Client_Time"] = associator_Register_Client_Time
				associator_json["Register_Client_MechineTime"] = associator_Register_Client_MechineTime
				associator_json["Receipt_Print_Manager"] = associator_Receipt_Print_Manager
				associator_list_table = append(associator_list_table, associator_json)
				fmt.Println("sqlite3数据库当前获取会员列表到的会员号" + associator_json["Name"].(string))
			}

		}
		query.Close()
		stmt.Close()

	}
	Sqlite_operation_lock.Unlock()

	//	var associator_json_text []byte
	//	associator_list_json := make(map[string]interface{})
	//	associator_list_json["Associator_List"] = associator_list_table
	//	associator_json_text, err = json.Marshal(associator_list_json)
	//	fmt.Println("sqlite3数据库获取会员列表" + string(associator_json_text))
	if len(associator_list_table) <= 0 {
		err = errors.New("no_such_associator")
	}
	return associator_list_table, err
}
