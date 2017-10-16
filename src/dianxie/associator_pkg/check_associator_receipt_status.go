package associator_pkg

import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Check_associator_receipt_status(associator Associator_s) (string, error) {
	Sqlite_operation_lock.Lock()
	var err error
	var res string
	sql_text := "SELECT Receipt_Print_Status,count(*) FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"
	var stmt *sql.Stmt
	stmt, err = Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库查询此会员打印收据状态sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库查询此会员打印收据状态sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		err = errors.New("sqlite3数据库查询此会员打印收据状态sql编译错误(err:" + err.Error() + ")")
		Sqlite_operation_lock.Unlock()
		return "", err

	}
	var query *sql.Rows
	query, err = stmt.Query(associator.Associator_Name, associator.Associator_Class, associator.Associator_Sex, associator.Associator_Phone_Number)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库查询此会员打印收据状态查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库查询此会员打印收据状态查询错误(err:" + err.Error() + ")")
		err = errors.New("sqlite3数据库查询此会员打印收据状态查询错误(err:" + err.Error() + ")")
	} else {
		var count int64
		query.Next()
		var receipt_print_status string
		query.Scan(&receipt_print_status, &count)
		query.Close()
		stmt.Close()

		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库查询此会员打印收据状态错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库查询此会员打印收据状态错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("sqlite3数据库查询此会员打印收据状态错误(err:此用户已存在且有为此信息的会员有一个以上)")
		} else if count == 1 {
			//此用户已存在
			mylog.ErrorLog("sqlite3数据库查询此会员打印收据状态成功！！\r\n")
			fmt.Println("sqlite3数据库查询此会员打印收据状态成功！！")
			res = receipt_print_status

		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库查询此会员打印收据状态错误(err:此用户不存在)\r\n")
			fmt.Println("sqlite3数据库查询查询此会员打印收据状态错误(err:此用户不存在)")
			err = errors.New("sqlite3数据库查询此会员打印收据状态错误(err:此用户不存在)")
		}
	}
	Sqlite_operation_lock.Unlock()
	return res, err
}
