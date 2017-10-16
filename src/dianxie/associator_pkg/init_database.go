package associator_pkg

import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Associator_s struct {
	Associator_Name                        string
	Associator_Number                      string
	Associator_Class                       string
	Associator_Sex                         string
	Associator_Phone_Number                string
	Associator_Card_Number                 string
	Associator_QQ_Number                   string
	Associator_Wechat_Number               string
	Associator_Register_Ip                 string
	Associator_Register_Mac                string
	Associator_Register_Server_Time        string
	Associator_Register_Server_MechineTime string
	Associator_Register_Client_Time        string
	Associator_Register_Client_MechineTime string
	Receipt_Print_Status                   string
	Receipt_Print_Ip                       string
	Receipt_Print_Mac                      string
	Receipt_Print_Server_Time              string
	Receipt_Print_Server_MechineName       string
	Receipt_Print_Client_Time              string
	Receipt_Print_Client_MechineName       string
	Receipt_Print_Manager                  string
}

var Db *sql.DB

//定义一个sqlite3的全局锁
var Sqlite_operation_lock *sync.Mutex

func Init_sqllite_databse() {

	var err error
	Db, err = sql.Open("sqlite3", appconf.Tf.SqlDb)
	if err != nil {
		mylog.ErrorLog("打开sqlite3数据库打开错误:文件不存在或者格式错误(err:" + err.Error() + ")\r\n")
		fmt.Println("打开sqlite3数据库打开错误:文件不存在或者格式错误(err:" + err.Error() + ")")
		panic(err)
	}

	_, err = Db.Query("SELECT * FROM " + appconf.Tf.SqlAssociatorTable)
	if err != nil {
		mylog.ErrorLog("打开sqlite3数据库数据表错误:数据表不存在或者格式错误(err:" + err.Error() + ")\r\n")
		fmt.Println("打开sqlite3数据库数据表错误:数据表不存在或者格式错误(err:" + err.Error() + ")")
		create_associator_table()
		//panic(err)
	}
	fmt.Println("sqlite3数据库初始化成功！！(addr：" + appconf.Tf.SqlDb + ")")
	Sqlite_operation_lock = new(sync.Mutex)
}
func create_associator_table() {

	var create_associator_table_sql_text string = `
	CREATE TABLE ` + appconf.Tf.SqlAssociatorTable + ` (
	    Associator_Name                         TEXT,
	    Associator_Number                       TEXT,
	    Associator_Class                        TEXT,
	    Associator_Sex                          TEXT,
	    Associator_Phone_Number                 TEXT,
	    Associator_Card_Number                  TEXT,
	    Associator_QQ_Number                    TEXT,
	    Associator_Wechat_Number                TEXT,
		Associator_Register_Ip                  TEXT,
	    Associator_Register_Mac                 TEXT,
	    Associator_Register_Server_Time         TEXT,
	    Associator_Register_Server_MechineTime  TEXT,
		Associator_Register_Client_Time         TEXT,
	    Associator_Register_Client_MechineTime  TEXT,
	    Receipt_Print_Status                    TEXT,
	    Receipt_Print_Ip                        TEXT,
	    Receipt_Print_Mac                       TEXT,
		Receipt_Print_Server_Time               TEXT,
	    Receipt_Print_Server_MechineName        TEXT,
	    Receipt_Print_Client_Time               TEXT,
	    Receipt_Print_Client_MechineTime        TEXT,
	    Receipt_Print_Manager                   TEXT
	)
	`
	var create_number_table_sql_text string = `
	CREATE TABLE ` + appconf.Tf.SqlAssociatorNumberTable + ` (
		Number_TYPE                     TEXT,
	    Associator_Number               TEXT
	)
	`
	var create_administrator_table_sql_text string = `
	CREATE TABLE ` + appconf.Tf.SqlAdministratorTable + ` (
		Administrator_Name                 TEXT,
	    Administrator_Passwd               TEXT
	)
	`
	var init_number_table_sql_text string = `
	INSERT INTO ` + appconf.Tf.SqlAssociatorNumberTable + ` (
		Number_TYPE,
		Associator_Number
	) values('max','0')
	`
	var err error
	mylog.ErrorLog("尝试修复sqlite3数据库数据库！！\r\n")
	fmt.Println("尝试修复sqlite3数据库！")
	_, err = Db.Exec(create_associator_table_sql_text) //这里使用Exec函数，因为这里是执行，经测试Query函数执行失败
	if err != nil {
		mylog.ErrorLog("尝试修复sqlite3数据库失败(err:" + err.Error() + ")\r\n")
		fmt.Println("尝试修复sqlite3数据库失败(err:" + err.Error() + ")")
		Db.Close()
		panic(err)
	}
	_, err = Db.Exec(create_number_table_sql_text) //这里使用Exec函数，因为这里是执行，经测试Query函数执行失败
	if err != nil {
		mylog.ErrorLog("尝试修复sqlite3数据库失败(err:" + err.Error() + ")\r\n")
		fmt.Println("尝试修复sqlite3数据库失败(err:" + err.Error() + ")")
		Db.Close()
		panic(err)
	}
	_, err = Db.Exec(create_administrator_table_sql_text) //这里使用Exec函数，因为这里是执行，经测试Query函数执行失败
	if err != nil {
		mylog.ErrorLog("尝试修复sqlite3数据库失败(err:" + err.Error() + ")\r\n")
		fmt.Println("尝试修复sqlite3数据库失败(err:" + err.Error() + ")")
		Db.Close()
		panic(err)
	}
	_, err = Db.Exec(init_number_table_sql_text) //这里使用Exec函数，因为这里是执行，经测试Query函数执行失败
	if err != nil {
		mylog.ErrorLog("尝试修复sqlite3数据库失败(err:" + err.Error() + ")\r\n")
		fmt.Println("尝试修复sqlite3数据库失败(err:" + err.Error() + ")")
		Db.Close()
		panic(err)
	}
	fmt.Println("尝试修复sqlite3数据库成功！！")
}
