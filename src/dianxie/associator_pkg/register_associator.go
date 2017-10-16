package associator_pkg

import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"
	"fmt"

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

func Register_associator(associator Associator_s) error {
	Sqlite_operation_lock.Lock()
	sql_text := "SELECT count(*) FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"

	stmt, err := Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库注册会员2sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库注册会员2sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		Sqlite_operation_lock.Unlock()
		return err

	}
	var query *sql.Rows
	query, err = stmt.Query(associator.Associator_Name, associator.Associator_Class, associator.Associator_Sex, associator.Associator_Phone_Number)
	if err != nil {
		//不存在此用户
		stmt.Close()
		query.Close()
		database_register_associator(associator)
	} else {
		var count int64
		query.Next()
		query.Scan(&count)
		stmt.Close()
		query.Close()
		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库注册会员时错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库注册会员时错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("严重错误:此用户已存在且有为此信息的会员有一个以上")
		} else if count == 1 {
			//此用户已存在
			mylog.ErrorLog("sqlite3数据库注册会员错误(err:此用户已存在)\r\n")
			fmt.Println("sqlite3数据库注册会员错误(err:此用户已存在)")
			err = errors.New("错误:此用户已存在")
		} else if count == 0 {
			//不存在此用户
			err = database_register_associator(associator)
		}
	}
	Sqlite_operation_lock.Unlock()
	return err
}
func database_register_associator(associator Associator_s) error {
	sql_text := "INSERT INTO " + appconf.Tf.SqlAssociatorTable
	//开始定义列名的
	sql_text += " (Associator_Name,"
	sql_text += " Associator_Number,"
	sql_text += " Associator_Class,"
	sql_text += " Associator_Sex,"
	sql_text += " Associator_Phone_Number,"
	sql_text += " Associator_Card_Number,"
	sql_text += " Associator_QQ_Number,"
	sql_text += " Associator_Wechat_Number,"
	sql_text += " Associator_Register_Ip,"
	sql_text += " Associator_Register_Mac,"
	sql_text += " Associator_Register_Server_Time,"
	sql_text += " Associator_Register_Server_MechineTime,"
	sql_text += " Associator_Register_Client_Time,"
	sql_text += " Associator_Register_Client_MechineTime,"
	sql_text += " Receipt_Print_Status"

	sql_text += " ) values("
	//开始填数据标识符
	sql_text += "?,?,?,?,?,?,?,?,?,?,?,?,?,?,?"
	sql_text += ")"

	st, err := Db.Begin()
	if err != nil {
		fmt.Println("sqlite3数据库注册会员状态失败(err:" + err.Error() + ")")
		return err
	}
	var stmt *sql.Stmt
	stmt, err = st.Prepare(sql_text)
	if err != nil {
		st.Rollback()
		mylog.ErrorLog("sqlite3数据库注册会员sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库注册会员时sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		return err

	}
	associator.Receipt_Print_Status = "no_proceed"
	var result sql.Result
	result, err = stmt.Exec(
		associator.Associator_Name,
		associator.Associator_Number,
		associator.Associator_Class,
		associator.Associator_Sex,
		associator.Associator_Phone_Number,
		associator.Associator_Card_Number,
		associator.Associator_QQ_Number,
		associator.Associator_Wechat_Number,
		associator.Associator_Register_Ip,
		associator.Associator_Register_Mac,
		associator.Associator_Register_Server_Time,
		associator.Associator_Register_Server_MechineTime,
		associator.Associator_Register_Client_Time,
		associator.Associator_Register_Client_MechineTime,
		associator.Receipt_Print_Status)
	if err != nil {
		//插入错误
		stmt.Close()
		st.Rollback()

		fmt.Println("sqlite3数据库注册会员失败(err:" + err.Error() + ")")
	} else {
		var affect int64
		affect, err = result.RowsAffected()
		if affect == 1 {
			st.Commit()
			fmt.Println("sqlite3数据库注册会员成功(message:name:" + associator.Associator_Name + " class:" + associator.Associator_Class + " sex:" + associator.Associator_Sex + " Phone_Number:" + associator.Associator_Phone_Number + ")")
		} else {
			st.Rollback()
			fmt.Println("sqlite3数据库注册会员失败2(err:" + err.Error() + ")")
			err = errors.New("sqlite3数据库注册会员失败2(err:" + err.Error() + ")")
		}

	}
	return err
}
