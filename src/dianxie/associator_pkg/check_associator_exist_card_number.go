package associator_pkg

import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Check_associator_exist_card_id(associator Associator_s) (bool, error) {
	Sqlite_operation_lock.Lock()
	var err error
	var res bool = false
	sql_text := "SELECT Associator_Card_Number,count(*) FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"
	var stmt *sql.Stmt
	stmt, err = Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库查询此会员会卡号是否存在sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库查询此会员会卡号是否存在sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		err = errors.New("sqlite3数据库查询此会员会卡号是否存在sql编译错误(err:" + err.Error() + ")")
		Sqlite_operation_lock.Unlock()
		return false, err

	}
	var query *sql.Rows
	query, err = stmt.Query(associator.Associator_Name, associator.Associator_Class, associator.Associator_Sex, associator.Associator_Phone_Number)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库查询此会员会卡号是否存在查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库查询此会员会卡号是否存在查询错误(err:" + err.Error() + ")")
		err = errors.New("sqlite3数据库查询此会员会卡号是否存在查询错误(err:" + err.Error() + ")")
	} else {
		var count int64
		query.Next()
		var associator_card_number string
		query.Scan(&associator_card_number, &count)
		query.Close()
		stmt.Close()

		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库查询此会员会卡号是否存在错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库查询此会员会卡号是否存在错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("sqlite3数据库查询此会员会卡号是否存在错误(err:此用户已存在且有为此信息的会员有一个以上)")
		} else if count == 1 {
			//此用户已存在
			mylog.ErrorLog("sqlite3数据库查询此会员会卡号是否存在成功！！\r\n")
			fmt.Println("sqlite3数据库查询此会员会卡号是否存在成功！！")
			if associator_card_number != "" {
				res = true

			} else {
				res = false
			}

		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库查询此会员会卡号是否存在错误(err:此用户不存在)\r\n")
			fmt.Println("sqlite3数据库查询此会员会卡号是否存在错误(err:此用户不存在)")
			err = errors.New("sqlite3数据库查询此会员会卡号是否存在错误(err:此用户不存在)")
		}
	}
	Sqlite_operation_lock.Unlock()
	return res, err
}
