package associator_pkg

import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Login_administrator(name string, passwd string) (bool, error) {
	Sqlite_operation_lock.Lock()
	var err error
	var res bool = false
	sql_text := "SELECT Administrator_Name,count(*) FROM " + appconf.Tf.SqlAdministratorTable
	sql_text += " WHERE Administrator_Name=?"
	sql_text += " AND Administrator_Passwd=?"
	var stmt *sql.Stmt
	stmt, err = Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库管理员登录sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库管理员登录sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		err = errors.New("sqlite3数据库管理员登录sql编译错误(err:" + err.Error() + ")")
		Sqlite_operation_lock.Unlock()
		return false, err

	}
	var query *sql.Rows
	query, err = stmt.Query(name, passwd)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库管理员登录查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库管理员登录查询错误(err:" + err.Error() + ")")
		err = errors.New("sqlite3数据库管理员登录查询错误(err:" + err.Error() + ")")
	} else {
		var count int64
		query.Next()
		var administrator_Name string
		query.Scan(&administrator_Name, &count)
		query.Close()
		stmt.Close()

		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库管理员登录错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库管理员登录错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("sqlite3数据库管理员登录错误(err:此用户已存在且有为此信息的会员有一个以上)")
		} else if count == 1 {
			//此用户已存在
			mylog.ErrorLog("sqlite3数据库管理员登录成功！！\r\n")
			fmt.Println("sqlite3数据库查询管理员登录成功！！")
			if administrator_Name != "" {
				res = true

			} else {
				res = false
			}

		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库管理员登录错误(err:此用户不存在)\r\n")
			fmt.Println("sqlite3数据库管理员登录错误(err:此用户不存在)")
			err = errors.New("sqlite3数据库管理员登录错误(err:此用户不存在)")
		}
	}
	Sqlite_operation_lock.Unlock()
	return res, err
}
