package myhttp

import (
	"dianxie/appconf"
	"dianxie/mylog"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Http_server_init() {
	Server := http.Server{
		Addr:           appconf.Tf.HttpAddr,
		Handler:        &MyHandle{},
		ReadTimeout:    appconf.Tf.HttpReadTimeout,
		WriteTimeout:   appconf.Tf.HttpWriteTimeout,
		MaxHeaderBytes: appconf.Tf.HttpMaxHeaderBytes,
	}
	Reg_url()
	var err error
	if appconf.Tf.HttpsEnable == "false" {
		err = Server.ListenAndServe()
	} else if appconf.Tf.HttpsEnable == "true" {
		_, err = ioutil.ReadFile(appconf.Tf.HttpsCertFile)
		if err != nil {
			fmt.Println("ReadHttpsCertFile err:", err)
			panic("ReadHttpsCertFile err:" + err.Error())
		}
		_, err = ioutil.ReadFile(appconf.Tf.HttpsKeyFile)
		if err != nil {
			fmt.Println("ReadHttpsKeyFile err:", err)
			panic("ReadHttpsKeyFile err:" + err.Error())
		}
		err = Server.ListenAndServeTLS(appconf.Tf.HttpsCertFile, appconf.Tf.HttpsKeyFile)
	} else {
		panic("httpsserver init failed err:not know httpserver protocol")
	}

	if err != nil {
		mylog.ErrorLog("HTTP SERVER 启动失败！")
	} else {
		mylog.ErrorLog("HTTP SERVER 启动成功！")
	}
}
