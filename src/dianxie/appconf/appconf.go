package appconf

/*
*这是个软件的配置文件的解析包
 */
import (
	"dianxie/mylog"
	"goconf"
	"net/http"
	"strings"
	"time"
)

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

var Tf *ServerConfig

func Read_conf() { //读取并解析配置文件函数

	conf := goconf.New()
	Tf = &ServerConfig{}
	if err := conf.Parse("../conf/app.conf"); err != nil {
		mylog.ErrorLog("解析配置文件错误:文件不存在或者格式错误\r\n")
	} else {
		core := conf.Get("http")
		if core == nil {
			mylog.ErrorLog("解析配置文件错误:配置格式错误\r\n")
		} else {

			if err := conf.Unmarshal(Tf); err != nil {
				mylog.ErrorLog("解析配置文件错误:配置格式错误\r\n")
			}
		}
	}
	Tf.SqlDb = strings.Replace(Tf.SqlDb, " ", "", -1)
	if Tf.SqlDb == "" {
		Tf.SqlDb = "../database/Electronics_Association_Associator_Database.db"
	}
	Tf.SqlAssociatorTable = strings.Replace(Tf.SqlAssociatorTable, " ", "", -1)
	if Tf.SqlAssociatorTable == "" {
		Tf.SqlAssociatorTable = "Associator"
	}

	Tf.SqlAssociatorNumberTable = strings.Replace(Tf.SqlAssociatorNumberTable, " ", "", -1)
	if Tf.SqlAssociatorNumberTable == "" {
		Tf.SqlAssociatorNumberTable = "AssociatorNumber"
	}
	Tf.SqlAdministratorTable = strings.Replace(Tf.SqlAdministratorTable, " ", "", -1)
	if Tf.SqlAdministratorTable == "" {
		Tf.SqlAdministratorTable = "Administrator"
	}

	Tf.HttpAddr = strings.Replace(Tf.HttpAddr, " ", "", -1)
	if Tf.HttpAddr == "" {
		Tf.HttpAddr = "0.0.0.0:8686"
	}
	Tf.HttpsEnable = strings.Replace(Tf.HttpsEnable, " ", "", -1)
	if Tf.HttpsEnable == "" {
		Tf.HttpsEnable = "false"
	}
	Tf.HttpsCertFile = strings.Replace(Tf.HttpsCertFile, " ", "", -1)
	if Tf.HttpsCertFile == "" {
		Tf.HttpsCertFile = "../httpscertificate/HttpsCert.pem"
	}
	Tf.HttpsKeyFile = strings.Replace(Tf.HttpsKeyFile, " ", "", -1)
	if Tf.HttpsKeyFile == "" {
		Tf.HttpsKeyFile = "../httpscertificate/HttpsKey.pem"
	}
	if Tf.HttpPort <= 0 {
		Tf.HttpPort = 8686
	}
	if Tf.UdpServerPort <= 0 {
		Tf.UdpServerPort = 3000
	}
	if Tf.UdpClientPort <= 0 {
		Tf.UdpClientPort = 3000
	}
	if Tf.UdpHeartbeakTime <= 0 {
		Tf.UdpHeartbeakTime = 700
	}
	if Tf.HttpReadTimeout <= 0 {
		Tf.HttpReadTimeout = 6 * time.Second
	}
	if Tf.HttpWriteTimeout <= 0 {
		Tf.HttpWriteTimeout = 6 * time.Second
	}
	if Tf.HttpMaxHeaderBytes <= 0 {
		Tf.HttpMaxHeaderBytes = 3 * http.DefaultMaxHeaderBytes
	}

}
