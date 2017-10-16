package associator_pkg

/*
*函数名:修改会员卡相关函数
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
*函数名:修改会员卡的入口函数
*参数：associator 会员信息
 */
func Change_associator_card_id(associator Associator_s) error {
	Sqlite_operation_lock.Lock()
	var err error
	if associator.Associator_Card_Number == "" {
	} else {
		if check_this_card_id_is_exsite(associator) == true {
			mylog.ErrorLog("sqlite3数据库此会员卡号已存在！！\r\n")
			fmt.Println("sqlite3数据库此会员卡号已存在！！")
			Sqlite_operation_lock.Unlock()
			err = errors.New("sqlite3数据库此会员卡号已存在！！")
			return err
		}
	}

	sql_text := "SELECT Associator_Card_Number,count(*) FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"
	var stmt *sql.Stmt
	stmt, err = Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库修改会卡号sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库修改会卡号sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		err = errors.New("sqlite3数据库修改会卡号sql编译错误(err:" + err.Error() + ")")
		Sqlite_operation_lock.Unlock()
		return err

	}
	var query *sql.Rows
	query, err = stmt.Query(associator.Associator_Name, associator.Associator_Class, associator.Associator_Sex, associator.Associator_Phone_Number)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库修改会卡号查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库修改会卡号查询错误(err:" + err.Error() + ")")
		err = errors.New("sqlite3数据库修改会卡号查询错误(err:" + err.Error() + ")")
	} else {
		var count int64
		query.Next()
		var associator_card_number string
		query.Scan(&associator_card_number, &count)
		query.Close()
		stmt.Close()

		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库修改会卡号错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库修改会卡号错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("sqlite3数据库修改会卡号错误(err:此用户已存在且有为此信息的会员有一个以上)")
		} else if count == 1 {
			//此用户已存在
			if associator_card_number != associator.Associator_Card_Number {
				database_change_associator_card_id(associator)
				mylog.ErrorLog("sqlite3数据库会员卡号修改成功！！\r\n")
				fmt.Println("sqlite3数据库会员卡号修改成功！！")

			} else {
				if associator_card_number == "" {

					mylog.ErrorLog("sqlite3数据库会员卡号修改成功！！\r\n")
					fmt.Println("sqlite3数据库会员卡号修改成功！！")
				} else {
					mylog.ErrorLog("sqlite3数据库修改会卡号错误(err:此用户卡号已经属于用户)\r\n")
					fmt.Println("sqlite3数据库修改会卡号错误(err:此用户卡号已经属于用户)")
					err = errors.New("sqlite3数据库修改会卡号错误(err:此用户卡号已经属于用户)")
				}

			}

		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库修改会卡号错误(err:此用户不存在)\r\n")
			fmt.Println("sqlite3数据库修改会卡号错误(err:此用户不存在)")
			err = errors.New("sqlite3数据库修改会卡号错误(err:此用户不存在)")
		}
	}
	Sqlite_operation_lock.Unlock()
	return err
}

/*
*函数名:检查次会员卡是否已存在
*参数：associator 会员信息
 */
func check_this_card_id_is_exsite(associator Associator_s) bool {
	//查询此卡号是否已存在
	sql_text := "SELECT Associator_Card_Number,count(*) FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Associator_Card_Number=?"

	stmt, err := Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库查询此卡号是否已存在错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库查询此卡号是否已存在错误(err:" + err.Error() + ")")
		//panic(err)
		return false
	}
	var query *sql.Rows

	query, err = stmt.Query(associator.Associator_Card_Number)
	if err != nil {
		//查询错误，可能是不存在此用户
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库查询此卡号是否已存在2成功(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库查询此卡号是否已存在2成功(err:" + err.Error() + ")")
		return false
	} else {
		var count int64
		query.Next()
		var associator_card_number string
		query.Scan(&associator_card_number, &count)
		query.Close()
		stmt.Close()

		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库查询此卡号是否已存在错误(err:此卡号已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库查询此卡号是否已存在错误(err:此卡号已存在且有为此信息的会员有一个以上)")
			return false
		} else if count == 1 {
			//此用户已存在
			if associator_card_number != "" {
				mylog.ErrorLog("sqlite3数据库查询此卡号是否已存在成功\r\n")
				fmt.Println("sqlite3数据库查询此卡号是否已存在成功")
				return true
			} else {
				mylog.ErrorLog("sqlite3数据库查询此卡号是否已存在错误(err:此用户卡号不存在)\r\n")
				fmt.Println("sqlite3数据库查询此卡号是否已存在错误(err:此用户卡号不存在)")
				return false
			}

		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库查询此卡号是否已存在成功(err:此卡号不存在)\r\n")
			fmt.Println("sqlite3数据库查询此卡号是否已存在成功(err:此卡号不存在)")
			return false
		}
	}
	return false

}

/*
*函数名:修改会员卡的执行函数
*参数：associator 会员信息
 */
func database_change_associator_card_id(associator Associator_s) {
	sql_text := "UPDATE " + appconf.Tf.SqlAssociatorTable
	sql_text += " SET Associator_Card_Number=?"
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"
	st, err := Db.Begin()
	if err != nil {
		fmt.Println("sqlite3数据库修改会员卡号2失败(err:" + err.Error() + ")")
		return
	}
	var stmt *sql.Stmt
	stmt, err = st.Prepare(sql_text)
	if err != nil {
		st.Rollback()
		mylog.ErrorLog("sqlite3数据库修改会员卡号2sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库修改会员卡号2sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		return

	}
	var result sql.Result
	result, err = stmt.Exec(
		associator.Associator_Card_Number,
		associator.Associator_Name,
		associator.Associator_Class,
		associator.Associator_Sex,
		associator.Associator_Phone_Number)
	if err != nil {
		//更新错误
		stmt.Close()
		st.Rollback()

		fmt.Println(err.Error() + "sqlite3数据库修改会员卡号2失败(err:" + err.Error() + ")")
	} else {
		var affect int64
		affect, err = result.RowsAffected()
		stmt.Close()
		if affect == 1 {
			st.Commit()
			fmt.Println("sqlite3数据库修改会员卡号2成功(message:name:" + associator.Associator_Name + " class:" + associator.Associator_Class + " sex:" + associator.Associator_Sex + " Phone_Number:" + associator.Associator_Phone_Number + ")")
		} else {
			st.Rollback()
			fmt.Println(err.Error() + "sqlite3数据库修改会员卡号2失败(err:" + err.Error() + ")")
		}

	}
}
