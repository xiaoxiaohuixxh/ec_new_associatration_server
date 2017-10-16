package associator_pkg

import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Delete_associator(associator Associator_s) (string, error) {
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
		mylog.ErrorLog("sqlite3数据库删除此会员sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库删除此会员sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		err = errors.New("sqlite3数据库删除此会员sql编译错误(err:" + err.Error() + ")")
		Sqlite_operation_lock.Unlock()
		return "", err
	}
	var query *sql.Rows
	query, err = stmt.Query(associator.Associator_Name, associator.Associator_Class, associator.Associator_Sex, associator.Associator_Phone_Number)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库删除此会员查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库删除此会员查询错误(err:" + err.Error() + ")")
		err = errors.New("sqlite3数据库删除此会员查询错误(err:" + err.Error() + ")")
	} else {
		var count int64
		query.Next()
		var receipt_print_status string
		query.Scan(&receipt_print_status, &count)
		query.Close()
		stmt.Close()

		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库删除此会员错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库删除此会员错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("sqlite3数据库删除此会员错误(err:此用户已存在且有为此信息的会员有一个以上)")
		} else if count == 1 {
			//此用户已存在
			mylog.ErrorLog("sqlite3数据库删除此会员已存在！！\r\n")
			fmt.Println("sqlite3数据库删除此会员已存在！！")
			//res = receipt_print_status
			_, err = delete_associator_when_it_exist(associator)
			if err != nil {
				//出错。
				res = "sqlite3数据库删除此会员失败！！"
				err = errors.New("sqlite3数据库删除此会员失败！！")
			} else {
				//删除成功
				res = "sqlite3数据库删除此会员成功！！"

			}

		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库删除此会员错误(err:此用户不存在)\r\n")
			fmt.Println("sqlite3数据库删除此会员错误(err:此用户不存在)")
			err = errors.New("sqlite3数据库删除此会员错误(err:此用户不存在)")
		}
	}
	Sqlite_operation_lock.Unlock()
	return res, err
}
func delete_associator_when_it_exist(associator Associator_s) (string, error) {
	//存在此会员，删除此会员
	//开始删除此员
	delect_sql_text := "DELETE FROM " + appconf.Tf.SqlAssociatorTable
	delect_sql_text += " WHERE Associator_Name=?"
	delect_sql_text += " AND Associator_Class=?"
	delect_sql_text += " AND Associator_Sex=?"
	delect_sql_text += " AND Associator_Phone_Number=?"
	stmt, err := Db.Prepare(delect_sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库删除此会员sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库删除此会员sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		return "", err
	}
	var result sql.Result
	result, err = stmt.Exec(associator.Associator_Name, associator.Associator_Class, associator.Associator_Sex, associator.Associator_Phone_Number)
	if err != nil {
		//删除出现错误
		stmt.Close()
		fmt.Println("sqlite3数据库删除此会员失败(err:" + err.Error() + ")")
		err = errors.New("sqlite3数据库删除此会员失败(err:" + err.Error() + ")")
	} else {
		var affect int64
		affect, err = result.RowsAffected()
		stmt.Close()
		if affect == 1 {
			fmt.Println("sqlite3数据库删除此会员成功")
		} else {
			fmt.Println("sqlite3数据库删除此会员失败(err:" + err.Error() + ")")
			err = errors.New("sqlite3数据库删除此会员失败(err:" + err.Error() + ")")
		}

	}
	return "", err
}
