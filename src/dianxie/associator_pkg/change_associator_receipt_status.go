package associator_pkg

//20170728 写好了修改收据打印状态的基础函数！！
//20170729 据了解sqlite3有三种模式，有可能不加全局锁,为了防止冲突我写的时候加全局锁！以后了解清楚再把锁去掉
//20170729 28号的时候忘记在插入的时候加上回滚操作了！！现在加上！
import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Change_associator_receipt_status(associator Associator_s) error {
	Sqlite_operation_lock.Lock()
	sql_text := "SELECT count(*) FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"

	stmt, err := Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库修改会员收据状态sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库修改会员收据状态sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		Sqlite_operation_lock.Unlock()
		return err

	}
	var query *sql.Rows
	query, err = stmt.Query(associator.Associator_Name, associator.Associator_Class, associator.Associator_Sex, associator.Associator_Phone_Number)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库修改会员收据状态查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库修改会员收据状态查询错误(err:" + err.Error() + ")")
	} else {
		var count int64
		query.Next()
		query.Scan(&count)
		query.Close()
		stmt.Close()

		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库修改会员收据状态错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库修改会员收据状态错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("sqlite3数据库修改会员收据状态错误(err:此用户已存在且有为此信息的会员有一个以上)")
		} else if count == 1 {
			//此用户已存在
			err = database_change_associator_receipt_status(associator)
		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库修改会员收据状态错误(err:此用户不存在)\r\n")
			fmt.Println("sqlite3数据库修改会员收据状态错误(err:此用户不存在)")
			err = errors.New("sqlite3数据库修改会员收据状态错误(err:此用户不存在)")
		}
	}
	Sqlite_operation_lock.Unlock()
	return err
}
func database_change_associator_receipt_status(associator Associator_s) error {
	sql_text := "UPDATE " + appconf.Tf.SqlAssociatorTable
	sql_text += " SET Receipt_Print_Status=?,"
	sql_text += " Receipt_Print_Ip=?,"
	sql_text += " Receipt_Print_Mac=?,"
	sql_text += " Receipt_Print_Server_Time=?,"
	sql_text += " Receipt_Print_Server_MechineName=?,"
	sql_text += " Receipt_Print_Client_Time=?,"
	sql_text += " Receipt_Print_Client_MechineTime=?,"
	sql_text += " Receipt_Print_Manager=?"
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"
	st, err := Db.Begin()
	if err != nil {
		fmt.Println("sqlite3数据库修改会员收据状态失败(err:" + err.Error() + ")")
		return err
	}
	var stmt *sql.Stmt
	stmt, err = st.Prepare(sql_text)
	if err != nil {
		st.Rollback()
		mylog.ErrorLog("sqlite3数据库修改会员收据状态2sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库修改会员收据状态2sql编译错误(err:" + err.Error() + ")")
		err = errors.New("sqlite3数据库修改会员收据状态2sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		return err

	}
	var result sql.Result
	result, err = stmt.Exec(
		associator.Receipt_Print_Status,
		associator.Receipt_Print_Ip,
		associator.Receipt_Print_Mac,
		associator.Receipt_Print_Server_Time,
		associator.Receipt_Print_Server_MechineName,
		associator.Receipt_Print_Client_Time,
		associator.Receipt_Print_Client_MechineName,
		associator.Receipt_Print_Manager,
		associator.Associator_Name,
		associator.Associator_Class,
		associator.Associator_Sex,
		associator.Associator_Phone_Number)
	if err != nil {
		//更新错误
		stmt.Close()
		st.Rollback()

		fmt.Println("sqlite3数据库修改会员收据状态失败(err:" + err.Error() + ")")
		err = errors.New("sqlite3数据库修改会员收据状态失败(err:" + err.Error() + ")")
	} else {
		var affect int64
		affect, err = result.RowsAffected()
		stmt.Close()
		if affect == 1 {
			st.Commit()
			fmt.Println("sqlite3数据库修改会员收据状态成功(message:name:" + associator.Associator_Name + " class:" + associator.Associator_Class + " sex:" + associator.Associator_Sex + " Phone_Number:" + associator.Associator_Phone_Number + ")")
		} else {
			st.Rollback()
			fmt.Println("sqlite3数据库修改会员收据状态失败(err:" + err.Error() + ")")
			err = errors.New("sqlite3数据库修改会员收据状态失败(err:" + err.Error() + ")")
		}

	}
	return err
}
