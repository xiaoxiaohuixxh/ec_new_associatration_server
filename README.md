---
date: 2017-10-14 21:42
status: public
title: 会员招新服务端说明
---

# 会员招新服务端说明 by xiaohui & taiqin
## 0.0 服务端源码说明
### 0.1 src目录结构说明
#### .
#### ├── dianxie　　　　　　　　　　　　　　代码目录
#### │   ├── appconf　　　　　　　 　　　　　软件配置的代码　　　
#### │   ├── associator_pkg　　　　　　　　　会员数据库的相关函数
#### │   ├── conf　　　　　　　　　　　　　　　存放软件的配置文件
#### │   ├── database　　　　　　　　　　　　　存放数据库
#### │   ├── httpscertificate　　　　　　　　　存放ｈｔｔｐｓ加密需要用到的证书
#### │   ├── log　　　　　　　　　　　　　　　　存放软件的日志
#### │   ├── myhttp　　　　　　　　　　　　　　存放ｈｔｔｐ相关的代码
#### │   ├── mylog　　　　　　　　　　　　　　　软件生成日志的相关代码
#### │   ├── sent_udp_heartbeak　　　　　　　　存放软件发送心跳包的相关代码
#### │   ├── server　　　　　　　　　　　　　　ｓｅｒｖｅｒ的主函数（入口文件）
#### │   └── soft_version　　　　　　　　　　　存放声明软件版本的代码
#### ├── github.com　　　　　　　　　　　　　　　从ｇｉｔｈｕｂ．ｃｏｍ下载的相关库
#### ├── goconf　　　　　　　　　　　　　　　　　存放软件处理配置文件的相关代码
#### ├── golang.org　　　　　　　　　　　　　　从ｇｏｌａｎｇ．ｏｒｇ下载的库
#### .
### 0.2 源码分析
#### 0.2.0 服务端主函数(入口函数)
##### 地址：/src/dianxie/server/main.go
##### 函数：main()服务端的主函数
``` 
    func main() {
        mylog.StartLog()//开始记录日志
        *************
        向日志打印启动信息
        *************
        appconf.Read_conf()//读取配置文件
        *************
        打印信息
        *************
        //数据库初始化
        associator_pkg.Init_sqllite_databse()
        
        W.Add(1)//添加一个信号量。
        
        //Http服务初始化
        go myhttp.Http_server_init()
        //Udp心跳广播服务初始化
        go sent_udp_heartbeak.Boardcast_udp_heartbeak_to_all_netcard()
        //for {
        //}
        W.Wait()//等待信号量，这里主要用于阻塞main函数，防止程序关闭，和while(1);作用类似，但是while(1);消耗cpu。
        associator_pkg.Db.Close()
    }
```
##### .
#### 0.2.1 软件配置的函数
#### 地址：/src/dianxie/appconf/appconf.go
##### ServerConfig结构体定义了配置文件的结构
```
type ServerConfig struct {
        	SqlDb                    string        `goconf:"database:SqlDb"`                    //数据库路径
        	SqlAssociatorTable       string        `goconf:"database:SqlAssociatorTable"`       //数据库的会员表的名字
        	SqlAssociatorNumberTable string        `goconf:"database:SqlAssociatorNumberTable"` //数据库的会员号表的名字
        	SqlAdministratorTable    string        `goconf:"database:SqlAdministratorTable"`    //数据库的管理员表的名字
        	HttpAddr                 string        `goconf:"http:HttpAddrAndPort"`              //http的地址和端口
        	HttpsEnable              string        `goconf:"http:HttpsEnable"`                  //使能https
        	HttpsCertFile            string        `goconf:"http:HttpsCertFile"`                //https的Certfile路径
        	HttpsKeyFile             string        `goconf:"http:HttpsKeyFile"`                 //https的keyfile的路径
        	HttpPort                 int           `goconf:"http:HttpOnlyPort"`                 //心跳包里包含的http的端口，必须跟上面定义的一样
        	HttpReadTimeout          time.Duration `goconf:"http:HttpReadTimeout"`              //http读超时
        	HttpWriteTimeout         time.Duration `goconf:"http:HttpWriteTimeout"`             //http写超时
        	HttpMaxHeaderBytes       int           `goconf:"http:HttpMaxHeaderBytes"`           //http读取最大字节
        	UdpServerPort            int           `goconf:"udp:UdpServerPort"`                 //udp服务端端口
        	UdpClientPort            int           `goconf:"udp:UdpClientPort"`                 //udp客户端端口
        	UdpHeartbeakTime         time.Duration `goconf:"udp:UdpHeartbeakTime"`              //udp心跳包刷新时间
        }
```
##### 解析配置文件函数
```

var Tf *ServerConfig //定义一个全局的配置结构体

func Read_conf() { //读取并解析配置文件函数
	conf := goconf.New()
	Tf = &ServerConfig{}
	if err := conf.Parse("../conf/app.conf"); err != nil {
    ************
    省略部分代码
    ************
	if err := conf.Unmarshal(Tf); err != nil {
    ************
    省略部分代码
    ************
**************
对读取到的配置信息进行合法性过滤和修正
**************
}
```
#### 0.2.2 数据库处理函数
#### 地址：src/dianxie/associator_pkg
##### init_database.go 数据库初始化
```
type Associator_s struct {//定义一个会员数据的结构体
	***********************
	省略会员数据的结构体
	***********************
}

var Db *sql.DB //声明一个全局的数据库连接句柄，用于初始化数据库后操作数据库不需要再次初始化数据库。
var Sqlite_operation_lock *sync.Mutex//定义一个sqlite3的全局锁，用于强制数据库操作串行执行，防止出现冲突。这里以后可以优化成并行操作。加快效率。
func Init_sqllite_databse() {
	var err error
	Db, err = sql.Open("sqlite3", appconf.Tf.SqlDb)//打开数据库文件
    **************
    省略部分代码
    **************
	_, err = Db.Query("SELECT * FROM " + appconf.Tf.SqlAssociatorTable)//检查数据库是否合法
    **************
    省略部分代码
    **************    
	fmt.Println("sqlite3数据库初始化成功！！(addr：" + appconf.Tf.SqlDb + ")")
	Sqlite_operation_lock = new(sync.Mutex)//初始化数据库串行锁
}
func create_associator_table() {

	var create_associator_table_sql_text string = `
	************
	创建会员表的sql语句
	***********
	var create_number_table_sql_text string = `
	************
	创建缓存会员号表的sql语句
	***********
	var create_administrator_table_sql_text string = `
	************
	创建管理员表的sql语句
	***********
	var init_number_table_sql_text string = `
	************
	创建缓存会员号表表的sql语句
	***********
    ************
	省略部分代码
	***********
	_, err = Db.Exec(create_associator_table_sql_text) //这里使用Exec函数，因为这里是执行，经测试Query函数执行失败
	************
	初始化失败就报错并关闭数据库连接
	***********
	_, err = Db.Exec(create_number_table_sql_text) //这里使用Exec函数，因为这里是执行，经测试Query函数执行失败
	************
	初始化失败就报错并关闭数据库连接
	***********
	_, err = Db.Exec(create_administrator_table_sql_text) //这里使用Exec函数，因为这里是执行，经测试Query函数执行失败
	************
	初始化失败就报错并关闭数据库连接
	***********
	_, err = Db.Exec(init_number_table_sql_text) //这里使用Exec函数，因为这里是执行，经测试Query函数执行失败
	************
	初始化失败就报错并关闭数据库连接
	***********
	fmt.Println("尝试修复sqlite3数据库成功！！")
}
```

##### cancel_associator_number.go 取消会员号的操作
```
func Cancel_associator_number(associator Associator_s) (string, error)//取消会员号的入口函数(注意入口函数都是全局的)，判断是否要执行取消会员号操作
func database_cancel_associator_number(associator Associator_s, associator_number string) (string, error)//执行取消会员号操作
```
##### change_associator_card_id.go 修改会员卡号
```
func Change_associator_card_id(associator Associator_s) error//修改会员卡号的入口函数(注意入口函数都是全局的)，判断是否要执行修改会员卡号操作
func check_this_card_id_is_exsite(associator Associator_s)//检查会员卡是否已存在
func database_change_associator_card_id(associator Associator_s)//执行修改会员卡号操作
```
##### change_associator_receipt_status.go 修改会员收据打印状态
```
func Change_associator_receipt_status(associator Associator_s) error//修改会员收据打印状态的入口函数(注意入口函数都是全局的)，判断是否要执行修改会员收据打印状态操作
func database_change_associator_receipt_status(associator Associator_s) error//修改会员收据打印状态的执行函数

```
##### check_associator_exist_card_number.go 检查会员是否存在会员卡号
```
func Check_associator_exist_card_id(associator Associator_s) (bool, error)执行检查会员是否存在会员卡号操作
```
##### check_associator_receipt_adminstrator.go 查询打印这个收据管理员信息
```
Check_associator_receipt_adminstrator(associator Associator_s) (string, error)//查询打印这个收据管理员信息
```
##### check_associator_receipt_status.go 检查会员收据打印状态，用于会返回会员的收据打印状态
```
Check_associator_receipt_status(associator Associator_s) (string, error)//检查会员收据打印状态，用于会返回会员的收据打印状态
```
##### delete_associator.go 删除会员
```
func Delete_associator(associator Associator_s) (string, error)//删除会员的入口函数(注意入口函数都是全局的)，判断是否要执行删除会员操作
func delete_associator_when_it_exist(associator Associator_s) (string, error)//执行删除会员操作
```
##### get_associator_list.go 获取会员列表
```
func Get_associator_list_by_receipt_status(receipt_Print_Status string) ([]interface{}, error//根据会员收据打印收据状态获取会员列表
```
##### get_associator_number.go 申请会员号
```
func Get_associator_number(associator Associator_s) (string, error)//申请会员号的入口函数(注意入口函数都是全局的)，判断如何进行会员号申请
func database_get_associator_number(associator Associator_s) (string, error)//当这个会员还没有会员号，这个函数会被Get_associator_number这个函数所调用。进行会员号申请操作。
func database_get_new_associator_number(st *sql.Tx) (string, error) //这个函数会被database_get_associator_number函数调用，进行从会员号缓存表里获取一个合法可用的会员号。
func database_get_new_associator_number_when_had_canceled(st *sql.Tx, stmt *sql.Stmt, query *sql.Rows) (string, error)//当缓存会员号表里有冗余的会员号这个函数会被database_get_new_associator_number函数调用并返回一个可用的会员号。
func database_get_new_associator_number_when_hadnt_canceled(st *sql.Tx, stmt *sql.Stmt, query *sql.Rows) (string, error)//当缓存会员号表里没有冗余的会员号这个函数会被database_get_new_associator_number函数调用并返回一个可用的会员号。
```
##### login_administrator.go 检查收据管理员是否合法，用于管理员登录
```
func Login_administrator(name string, passwd string) (bool, error) //检查收据管理员是否合法，合法的话第一个返回值会是true，否则false。第二个返回值是报错信息
```
##### register_associator.go 会员注册
```
func Register_associator(associator Associator_s) error//会员注册的入口函数(注意入口函数都是全局的)，判断是否进行会员注册。
func database_register_associator(associator Associator_s) error //执行会员注册
```
#### 0.2.3 http处理函数
#### 地址：src/dianxie/myhttp
##### myhttp.go http初始化函数
```
func Http_server_init()//初始化http，并且把MyHandle结构体绑定为请求回调接口。
```
##### analysis.go http路径解析函数
```
func (*MyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request)//当http数据请求来临时会被调用，w是你写数据的句柄，r是读请求数据的句柄。
```
##### cancel_associator_number_func.go 处理取消会员号请求的函数
```
func httpfun_cancel_associator_number(w http.ResponseWriter, r *http.Request)
```
##### change_associator_card_id_func.go 处理修改会员卡号请求的函数
```
func httpfun_change_associator_card_id(w http.ResponseWriter, r *http.Request)
```
##### change_associator_receipt_status_func.go 处理修改会员收据打印状态请求的函数
```
func httpfun_change_associator_receipt_status(w http.ResponseWriter, r *http.Request)
```
##### check_associator_exist_card_number_func.go 处理查询会员卡号请求的函数
```
func httpfun_check_associator_exist_card_number(w http.ResponseWriter, r *http.Request)
```
##### check_associator_receipt_adminstartor_func.go 处理查询处理此收据的管理员信息请求的函数
```
func httpfun_check_associator_receipt_adminstrator(w http.ResponseWriter, r *http.Request)
```
##### check_associator_receipt_status_func.go 处理查询此会员收据的打印状态请求的函数
```
func httpfun_check_associator_receipt_status(w http.ResponseWriter, r *http.Request)
```
##### delete_associator_func.go 处理删除会员请求的函数
```
func httpfun_delete_associator(w http.ResponseWriter, r *http.Request) 
```
##### get_associator_list_func.go 处理获取会员列表请求的函数
```
func httpfun_get_associator_list(w http.ResponseWriter, r *http.Request)
```
##### get_associator_number_func.go 处理申请会员号请求的函数
```
func httpfun_get_associator_number(w http.ResponseWriter, r *http.Request)
```
##### login_administrator_func.go 处理管理员登录请求的函数
```
func httpfun_login_administrator(w http.ResponseWriter, r *http.Request)
```
##### login_administrator_func.go 处理管理员登录请求的函数
```
func httpfun_login_administrator(w http.ResponseWriter, r *http.Request)
```
##### register_associator_func.go 处理会员注册请求的函数
```
func httpfun_register_associator(w http.ResponseWriter, r *http.Request) 
```
#### 0.2.4 软件日志处理函数
#### 地址：src/dianxie/mylog/mylog.go
##### 初始化日志处理
```
func StartLog()
```  
##### 打印错误日志信息
```
func ErrorLog(text string)
``` 
#### 0.2.5 软件的心跳包封包和发送函数
#### 地址：src/dianxie/sent_udp_heartbeak/sent_udp_boardcast.go
##### 初始化心跳包发送定时器
```
func Boardcast_udp_heartbeak_to_all_netcard()
```
##### 向本机的所有网卡发送心跳包，此函数会被定时器回调（支持ipv4和ipv6）
```
func sent_udp_heartbeak_to_all_netcard()
```
##### 向指定ip的网卡发送心跳包
```
func sent_udp_heartbeak(ip string)
```
##### 向指定ip的网卡发送心跳包
```
func sent_udp_heartbeak(ip string)
```
##### 向指定ip的网卡发送数据包，此函数会被sent_udp_heartbeak调用。
```
func sent_udp_boardcast_page(sent_addr net.UDPAddr, page []byte)
```
#### 0.2.0 服务端软件版本
##### 地址：/src/dianxie/soft_version/soft_version.go
```
var Soft_version string = "1.1.0"
```

#### .
## 1.0 服务端使用udp广播的方式来发送心跳包，心跳包里包含了服务端的ip,mac,版本等信息。封装心跳包的代码位于服务端源码的(src/dianxie/sent_udp_heartbeak/sent_udp_boardcast.go的func sent_udp_heartbeak(ip string)函数，参数ip是要发送心跳包的网卡的ip)
## 1.1 心跳包使用json结构封装。举例：
``` func sent_udp_heartbeak(ip string) {
        var heartbeak_json_text []byte
        heartbeak_table := make(map[string]interface{})
        heartbeak_table["Soft_version"] = soft_version.Soft_version
        heartbeak_table["Ip"] = ip
        heartbeak_table["HttpSEnable"] = appconf.Tf.HttpsEnable
        heartbeak_table["HttpServerPort"] = appconf.Tf.HttpPort
        heartbeak_json := make(map[string]interface{})
        heartbeak_json["Udp_Heartbeak"] = heartbeak_table
        heartbeak_json_text, _ = json.Marshal(heartbeak_json)
        var sent_addr net.UDPAddr
        sent_addr.IP = net.ParseIP(ip)
        sent_addr.Port = appconf.Tf.UdpServerPort
        sent_udp_boardcast_page(sent_addr, heartbeak_json_text)
        //fmt.Println(ip + "  " + string(heartbeak_json_text))
        }
        ```
### 心跳包例子：
```　
{"Udp_Heartbeak":{"HttpSEnable":"true","HttpServerPort":8686,"Ip":"192.168.1.182","Soft_version":"1.1.0"}} 
```
### Udp_Heartbeak说明后面的是心跳包的内容
### HttpSEnable 如果这个属性的值为true代表服务端开启了https，客户端必须使用https来传输数据。false反之。
### HttpServerPort 这个属性是服务端http/https服务所使用的端口号。
### Ip 这个属性是服务端的ip。
### Soft_version 这个属性是服务端的版本号。
## 2.0 服务端使用http或者https的方式与客户端进行数据交换。(严重推荐使用https)
### 2.1 服务端的http代码位于(src/dianxie/myhttp文件夹)
#### myhttp文件夹的文件名称说明:
#### 文件名包含_func的都是存放实现这一功能的函数的go文件
####
### 2.2 http的api说明。
#### 会员注册
##### 地址：/PC_APP_API/Register_Associator
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Register_Associator":{"Name":"会员名字","Class":"班级","Sex":"man或者girl","Phone_Number":"手机号码","QQ_Number":"QQ号码","Wechat_Number":"微信号码",,"Register_Time":"注册时间","Register_MechineTime":"注册时间戳","Register_Mac":"注册的机器的mac地址"}} 
```
#### 根据收据打印状态获取会员列表
##### 地址：/PC_APP_API/Get_Associator_List
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Get_Associator_List":{"Receipt_Print_Status":"no_proceed"}} 
```
#### //备注
#### Receipt_Print_Status = "no_proceed"   => 获取未进行中列表
#### Receipt_Print_Status = "proceed"         => 获取进行中列表
#### 申请会员号
##### 地址：/PC_APP_API/Get_Associator_Number
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Get_Associator_Number":{"Name":"姓名","Class":"班级","Sex":"性别","Phone_Number":"手机号码"}}
```
#### 修改会员的收据打印状态
##### 地址：/PC_APP_API/Change_Associator_Receipt_Status
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Change_Associator_Receipt_Status":{"Name":"姓名","Class":"班级","Sex":"性别","Phone_Number":"手机号码","Receipt_Print_Status":"no_proceed","Print_Time":"打印现行时间","Print_MechineTime":"打印现行时间戳","Print_Mac":"本地mac地址","Print_Manager":"管理员名字"}}
```
#### 修改会员的会员卡号
##### 地址：/PC_APP_API/Change_Associator_Card_Id
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Chang_Associator_Card_Number":{"Name":"姓名","Class":"班级","Sex":"性别","Phone_Number":"手机号码","Card_Number":"会员卡号"}}
```
#### 取消会员已申请的会员号
##### 地址：/PC_APP_API/Cancel_Associator_Number
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Cancel_Associator_Number":{"Name":"姓名","Class":"班级","Sex":"性别","Phone_Number":"手机号码","Number":"会员号码"}}
```
#### 查询此会员是否存在会员卡
##### 地址：/PC_APP_API/Check_Associator_Exist_Card_Number
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
```
{"Check_Associator_Exist_Card_Number":{"Name":"姓名","Class":"班级","Sex":"性别","Phone_Number":"手机号码"}}
```
#### 管理员登录
##### 地址：/PC_APP_API/Login_Administrator
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Login_Administrator":{"Name":"管理员账号","Passwd":"管理员密码"}}
```
#### 查询此会员当前收据打印状态
##### 地址：/PC_APP_API/Check_Associator_Receipt_Status
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Check_Associator_Receipt_Status":{"Name":"姓名","Class":"班级","Sex":"性别","Phone_Number":"手机号码"}}
```
#### 查询此会员当前收据打印的管理员
##### 地址：/PC_APP_API/Check_Associator_Receipt_Adminstartor
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Check_Associator_Receipt_Adminstartor":{"Name":"姓名","Class":"班级","Sex":"性别","Phone_Number":"手机号码"}}
```
#### 删除会员
##### 地址：/PC_APP_API/Delete_Associator
##### 请求方式：POST
##### 参数名称：data
##### 数据格式：json
##### json例子：
```
{"Delete_Associator":{"Name":"姓名","Class":"班级","Sex":"性别","Phone_Number":"手机号码"}}
```
### 参与人员
#### jiamei mingming taiqin zibo zhaoyong jiaji xiaohui(负责打杂) 
### 鸣谢全体参与招新的电协人。