package associator_pkg

//20170728 写好了修改收据打印状态的基础函数！！
//20170729 据了解sqlite3有三种模式，有可能不加全局锁,为了防止冲突我写的时候加全局锁！以后了解清楚再把锁去掉
//20170729 28号的时候忘记在插入的时候加上回滚操作了！！现在加上！
//20170729 在数据库添加了会员表，为了存放会员号被退回的会员号！！
//20170729 完成了申请会员会员号的代码编写，还没有测试
//20170730 已测试并修复申请会员会员号的代码
//20170730 编写并测试和修复了取消会员会员号的代码
import (
	"database/sql"
	"dianxie/appconf"
	"dianxie/mylog"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func Get_associator_number(associator Associator_s) (string, error) {
	Sqlite_operation_lock.Lock()
	sql_text := "SELECT Associator_Number,count(*) FROM " + appconf.Tf.SqlAssociatorTable
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"

	stmt, err := Db.Prepare(sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库申请会员会员号sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库申请会员会员号sql编译错误(err:" + err.Error() + ")")
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
		mylog.ErrorLog("sqlite3数据库申请会员会员号查询错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库修申请会员会员号查询错误(err:" + err.Error() + ")")
	} else {
		var count int64
		query.Next()
		query.Scan(&associator_number, &count)
		stmt.Close()
		query.Close()
		if count > 1 {
			//此用户已存在且有为此信息的会员有一个以上
			mylog.ErrorLog("sqlite3数据库申请会员会员号错误(err:此用户已存在且有为此信息的会员有一个以上)\r\n")
			fmt.Println("sqlite3数据库申请会员会员号错误(err:此用户已存在且有为此信息的会员有一个以上)")
			err = errors.New("sqlite3数据库申请会员会员号错误(err:此用户已存在且有为此信息的会员有一个以上)")
		} else if count == 1 {
			//此用户已存在
			if associator_number != "" {
				//此用户已有会员号
				mylog.ErrorLog("sqlite3数据库申请会员会员号错误(err:此用户已有会员号)\r\n")
				fmt.Println("sqlite3数据库申请会员会员号错误(err:此用户已有会员号)")
				err = errors.New("此用户已有会员号")
			} else {
				associator_number, err = database_get_associator_number(associator)
			}

		} else if count == 0 {
			//不存在此用户
			mylog.ErrorLog("sqlite3数据库申请会员会员号错误(err:此用户不存在)\r\n")
			fmt.Println("sqlite3数据库申请会员会员号错误(err:此用户不存在)")
			err = errors.New("sqlite3数据库申请会员会员号错误(err:此用户不存在)")
		}
	}
	Sqlite_operation_lock.Unlock()
	return associator_number, err
}
func database_get_new_associator_number_when_had_canceled(st *sql.Tx, stmt *sql.Stmt, query *sql.Rows) (string, error) {
	var associator_number string
	var err error
	query.Next()
	query.Scan(&associator_number)
	stmt.Close()
	query.Close()
	if associator_number != "" {
		//存在被取消的会员号，优先使用被取消的会员号！！
		//开始删除此被取消的会员号
		delect_sql_text := "DELETE FROM " + appconf.Tf.SqlAssociatorNumberTable
		delect_sql_text += " WHERE Associator_Number=?"
		delect_sql_text += " AND Number_TYPE=?"
		stmt, err = st.Prepare(delect_sql_text)
		if err != nil {
			mylog.ErrorLog("sqlite3数据库获取最新会员会员号2sql编译错误(err:" + err.Error() + ")\r\n")
			fmt.Println("sqlite3数据库获取最新会员会员号2sql编译错误(err:" + err.Error() + ")")
			//panic(err)
			return "", err
		}
		var result sql.Result
		result, err = stmt.Exec(associator_number, "cancelled")
		if err != nil {
			//删除出现错误
			stmt.Close()

			fmt.Println("sqlite3数据库获取最新会员会员号失败(err:" + err.Error() + ")")
		} else {
			var affect int64
			affect, err = result.RowsAffected()
			stmt.Close()
			if affect == 1 {
				fmt.Println("sqlite3数据库获取最新会员会员号成功")
			} else {
				fmt.Println("sqlite3数据库获取最新会员会员号失败(err:" + err.Error() + ")")
				err = errors.New("sqlite3数据库获取最新会员会员号失败(err:" + err.Error() + ")")
			}

		}
	} else {
		//不存在被取消的会员号，从现在最大的会员号继续获取！！
		fmt.Println("sqlite3数据库不存在被取消的会员号(err:no err)")
		associator_number, err = database_get_new_associator_number_when_hadnt_canceled(st, stmt, query)

	}
	return associator_number, err
}

func database_get_new_associator_number_when_hadnt_canceled(st *sql.Tx, stmt *sql.Stmt, query *sql.Rows) (string, error) {
	var associator_number string
	var err error
	//当不存在被取消的会员号时
	select_sql_text := "SELECT Associator_Number"
	select_sql_text += " FROM " + appconf.Tf.SqlAssociatorNumberTable
	select_sql_text += " WHERE Number_TYPE=?"
	stmt.Close()
	stmt, err = st.Prepare(select_sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库获取最新会员会员号sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库获取最新会员会员号sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		return "", err

	}
	//query.Close() //关闭query
	//开始获取最新的会员号------------------------------------------------------------------------------
	query, err = stmt.Query("max") //查询现在最大的会员号！
	if err != nil {
		//查询现在最大的会员号失败，直接返回失败
		stmt.Close()
		mylog.ErrorLog("sqlite3数据库查询现在最大的会员号失败(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库查询现在最大的会员号失败(err:" + err.Error() + ")")
		//panic(err)
		return "", err
	} else {
		//查询现在最大的会员号成功
		query.Next()
		query.Scan(&associator_number)
		query.Close()
		stmt.Close()

		associator_number_i, error := strconv.Atoi(associator_number)
		if error != nil {
			fmt.Println(associator_number + "sqlite3数据库查询现在最大的会员号失败(err:字符串转换成整数失败)")
			return "", err
		}
		associator_number_i = associator_number_i + 1
		associator_number = strconv.Itoa(associator_number_i) //数字变成字符串

		//开始更新最新的会员号到数据库------------------------------------------------------------------------------
		update_sql_text := "UPDATE " + appconf.Tf.SqlAssociatorNumberTable
		update_sql_text += " SET Associator_Number=?"
		update_sql_text += " WHERE Number_TYPE=?"
		stmt, err = st.Prepare(update_sql_text)
		if err != nil {
			mylog.ErrorLog("sqlite3数据库获取最新会员会员号2sql编译错误(err:" + err.Error() + ")\r\n")
			fmt.Println("sqlite3数据库获取最新会员会员号2sql编译错误(err:" + err.Error() + ")")
			//panic(err)
			return "", err
		}
		var result sql.Result
		result, err = stmt.Exec(associator_number, "max")
		if err != nil {
			//更新错误
			stmt.Close()

			fmt.Println("sqlite3数据库获取最新会员会员号时更新会员号到数据库失败(err:" + err.Error() + ")")
		} else {
			var affect int64
			affect, err = result.RowsAffected()
			stmt.Close()
			if affect == 1 {
				fmt.Println("sqlite3数据库获取最新会员会员号成功")
			} else {
				fmt.Println("sqlite3数据库获取最新会员会员号失败(err:" + err.Error() + ")")
				err = errors.New("sqlite3数据库获取最新会员会员号失败(err:" + err.Error() + ")")
			}

		}
	}
	return associator_number, err
}
func database_get_new_associator_number(st *sql.Tx) (string, error) {
	select_sql_text := "SELECT Associator_Number"
	select_sql_text += " FROM " + appconf.Tf.SqlAssociatorNumberTable
	select_sql_text += " WHERE Number_TYPE=?"

	var stmt *sql.Stmt
	var associator_number string
	stmt, err := st.Prepare(select_sql_text)
	if err != nil {
		mylog.ErrorLog("sqlite3数据库获取最新会员会员号sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库获取最新会员会员号sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		return "", err

	}
	var query *sql.Rows
	query, err = stmt.Query("cancelled") //查询被取消的会员号！
	if err != nil {
		//不存在被取消的会员号，从现在最大的会员号继续获取！！
		fmt.Println("sqlite3数据库不存在被取消的会员号(err:" + err.Error() + ")")
		associator_number, err = database_get_new_associator_number_when_hadnt_canceled(st, stmt, query)
		//当不存在被取消的会员号，从现在最大的会员号继续获取！！返回最新会员号-----------------------------------------------------------------------------------------------------------------

	} else {
		//存在被取消的会员号，优先使用被取消的会员号！！
		associator_number, err = database_get_new_associator_number_when_had_canceled(st, stmt, query)
		//上面是存在被取消的会员号，优先使用被取消的会员号！！而且在删除此被取消的会员号

	}
	return associator_number, err
}
func database_get_associator_number(associator Associator_s) (string, error) {
	sql_text := "UPDATE " + appconf.Tf.SqlAssociatorTable
	sql_text += " SET Associator_Number=?"
	sql_text += " WHERE Associator_Name=?"
	sql_text += " AND Associator_Class=?"
	sql_text += " AND Associator_Sex=?"
	sql_text += " AND Associator_Phone_Number=?"

	st, err := Db.Begin()
	if err != nil {
		fmt.Println("sqlite3数据库申请会员会员号失败(err:" + err.Error() + ")")
		return "", err
	}
	associator.Associator_Number, err = database_get_new_associator_number(st) //从数据库获取最新的会员号！！返回错误时已经回滚事务
	if err != nil {
		st.Rollback()
		fmt.Println("sqlite3数据库申请会员会员号失败(err:" + err.Error() + ")")
		return "", err
	}
	var stmt *sql.Stmt
	stmt, err = st.Prepare(sql_text)
	if err != nil {
		st.Rollback()
		mylog.ErrorLog("sqlite3数据库申请会员会员号2sql编译错误(err:" + err.Error() + ")\r\n")
		fmt.Println("sqlite3数据库申请会员会员号2sql编译错误(err:" + err.Error() + ")")
		//panic(err)
		return "", err

	}
	var result sql.Result
	result, err = stmt.Exec(
		associator.Associator_Number,
		associator.Associator_Name,
		associator.Associator_Class,
		associator.Associator_Sex,
		associator.Associator_Phone_Number)
	if err != nil {
		//更新错误
		stmt.Close()
		st.Rollback()

		fmt.Println("sqlite3数据库申请会员会员号失败(err:" + err.Error() + ")")
	} else {
		var affect int64
		affect, err = result.RowsAffected()
		stmt.Close()
		if affect == 1 {
			st.Commit()
			fmt.Println("sqlite3数据库申请会员会员号成功(message:name:" + associator.Associator_Name + " class:" + associator.Associator_Class + " sex:" + associator.Associator_Sex + " Phone_Number:" + associator.Associator_Phone_Number + ")")
		} else {
			st.Rollback()
			fmt.Println("sqlite3数据库申请会员会员号失败(err:" + err.Error() + ")")
			err = errors.New("sqlite3数据库申请会员会员号失败(err:" + err.Error() + ")")
		}

	}
	return associator.Associator_Number, err
}
