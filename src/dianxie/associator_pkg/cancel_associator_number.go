package associator_pkg

/*
*取消会员号相关函数
 */
import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

/*
*函数名:取消会员号入口函数
*参数：associator 会员信息
 */
func Cancel_associator_number(associator Associator_s) (string, error) {
	Sqlite_operation_lock.Lock()
	sql_text := "SELECT Associator_Number,count(*) FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"

	stmt, err := Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库取消会员会员号sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库取消会员会员号sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		Sqlite_operation_lock.Unlock()
		return "", err

	}
	var query *sql.Rows
	var associator_number string
	query, err = stmt.Query(associator.Associator_Name, associator.Associator_Class, associator.Associator_Sex, associator.Associator_Phone_Number)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库取消会员会员号查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库修取消会员会员号查询错误(err:" + err.Error() + ")")
	} else {
		var count int64
		query.Next()
		query.Scan(&associator_number, &count)
		stmt.Close()
		query.Close()
		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库取消会员会员号错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库取消会员会员号错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("sqlite3数据库取消会员会员号错误(err:此用户已存在且有为此信息的会员有一个以上)")
		} else if count == 1 {
			//此用户已存在
			if associator_number == "" {
				//此用户已有会员号
				mylog.ErrorLog("sqlite3数据库取消会员会员号错误(err:此用户没有会员号)\r\n")
				fmt.Println("sqlite3数据库取消会员会员号错误(err:此用户没有会员号)")
				err = errors.New("此用户没有会员号")
			} else {
				associator_number, err = database_cancel_associator_number(associator, associator_number)
			}

		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库取消会员会员号错误(err:此用户不存在)\r\n")
			fmt.Println("sqlite3数据库取消会员会员号错误(err:此用户不存在)")
			err = errors.New("sqlite3数据库取消会员会员号错误(err:此用户不存在)")
		}
	}
	Sqlite_operation_lock.Unlock()
	return associator_number, err
}

/*
*函数名:取消会员号的开始取消会员号函数
*参数：associator 会员信息
*     associator_number 会员号
 */
func database_cancel_associator_number(associator Associator_s, associator_number string) (string, error) {
	//开始取消会员号的数据库操作事务
	st, err := Db.Begin() //开始事务
	if err != nil {
		fmt.Println("sqlite3数据库取消会员会员号失败(err:" + err.Error() + ")")
		return associator_number, err
	}
	//把会员号写到会员号表取消表内
	insert_sql_text := "INSERT INTO " + appconf.Tf.SqlAssociatorNumberTable
	//开始定义列名的
	insert_sql_text += " (Associator_Number,"
	insert_sql_text += " Number_TYPE"

	insert_sql_text += " ) values("
	//开始填数据标识符
	insert_sql_text += "?,?"
	insert_sql_text += ")"

	stmt, err := st.Prepare(insert_sql_text)
	if err != nil {
		st.Rollback()
		mylog.ErrorLog("sqlite3数据库取消会员会员号2sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库取消会员会员号2sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		return associator_number, err

	}

	var result sql.Result
	result, err = stmt.Exec(associator_number, "cancelled")
	if err != nil {
		//插入错误
		stmt.Close()
		st.Rollback()

		fmt.Println("sqlite3数据库取消会员会员号失败(err:" + err.Error() + ")")

	} else {
		var affect int64
		affect, err = result.RowsAffected()

		if affect == 1 {
			//把会员号写到会员号表取消表内成功,把会员的会员号清空
			stmt.Close()
			update_sql_text := "UPDATE " + appconf.Tf.SqlAssociatorTable
			update_sql_text += " SET Associator_Number=?"
			update_sql_text += " WHERE Associator_Number=?"
			update_sql_text += " AND Associator_Name=?"
			update_sql_text += " AND Associator_Class=?"
			update_sql_text += " AND Associator_Sex=?"
			update_sql_text += " AND Associator_Phone_Number=?"
			stmt, err = st.Prepare(update_sql_text)
			if err != nil {
				st.Rollback()
				mylog.ErrorLog("sqlite3数据库取消会员会员号2sql编译错误(err:" + err.Error() + ")\r\n")
				fmt.Println("sqlite3数据库取消会员会员号2sql编译错误(err:" + err.Error() + ")")
				//panic(err)
				return associator_number, err

			}
			result, err = stmt.Exec(
				"",
				associator_number,
				associator.Associator_Name,
				associator.Associator_Class,
				associator.Associator_Sex,
				associator.Associator_Phone_Number)
			if err != nil {
				//插入错误
				stmt.Close()
				st.Rollback()

				fmt.Println("sqlite3数据库取消会员会员号失败(err:" + err.Error() + ")")
			} else {
				var affect int64
				affect, err = result.RowsAffected()
				stmt.Close()
				if affect == 1 {
					st.Commit()
					fmt.Println("sqlite3数据库取消会员会员号成功(message:name:" + associator.Associator_Name + " class:" + associator.Associator_Class + " sex:" + associator.Associator_Sex + " Phone_Number:" + associator.Associator_Phone_Number + ")")
				} else {
					st.Rollback()
					fmt.Println("sqlite3数据库取消会员会员号失败(err:" + err.Error() + ")")
					err = errors.New("sqlite3数据库取消会员会员号失败(err:" + err.Error() + ")")
				}
			}

		} else {
			st.Rollback()
			fmt.Println("sqlite3数据库取消会员会员号失败2(err:" + err.Error() + ")")
			err = errors.New("sqlite3数据库取消会员会员号失败2(err:" + err.Error() + ")")
		}

	}
	return associator_number, err
}
