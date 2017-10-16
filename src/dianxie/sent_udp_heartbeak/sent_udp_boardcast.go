package sent_udp_heartbeak

import (
	"dianxie/appconf"
	"dianxie/soft_version"
	"encoding/json"
	//"fmt"
	"net"
	"time"
)

//func sent_udp_page(sent_addr net.UDPAddr, rev_addr net.UDPAddr, page []byte) {
//	// 这里设置发送者的IP地址，自己查看一下自己的IP自行设定
//	//	laddr := net.UDPAddr{
//	//		IP:   net.IPv4(192, 168, 137, 224),
//	//		Port: 3000,
//	//	}
//	// 这里设置接收者的IP地址为广播地址
//	//	raddr := net.UDPAddr{
//	//		IP:   net.IPv4(255, 255, 255, 255),
//	//		Port: 3000,
//	//	}
//	conn, err := net.DialUDP("udp", &sent_ip, &rev_ip)
//	if err != nil {
//		println(err.Error())
//		return
//	}
//	conn.Write(page)
//	conn.Close()
//}
func sent_udp_boardcast_page(sent_addr net.UDPAddr, page []byte) {
	// 这里设置发送者的IP地址，自己查看一下自己的IP自行设定
	//	laddr := net.UDPAddr{
	//		IP:   net.IPv4(192, 168, 137, 224),
	//		Port: 3000,
	//	}
	// 这里设置接收者的IP地址为广播地址
	raddr := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: appconf.Tf.UdpClientPort,
	}
	conn, err := net.DialUDP("udp", &sent_addr, &raddr)
	if err != nil {
		//println(err.Error())
		return
	}
	conn.Write(page)
	conn.Close()
}
func sent_udp_heartbeak(ip string) {
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
func sent_udp_heartbeak_to_all_netcard() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		//fmt.Println(err)
		//os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				//fmt.Println(ipnet.IP.String())
				sent_udp_heartbeak(ipnet.IP.String())
			} else if ipnet.IP.To16() != nil {

				//fmt.Println(ipnet.IP.String())
				sent_udp_heartbeak(ipnet.IP.String())
			}

		}
	}

}
func Boardcast_udp_heartbeak_to_all_netcard() {
	timer1 := time.NewTicker(appconf.Tf.UdpHeartbeakTime * time.Millisecond)
	for {
		select {
		case <-timer1.C:
			sent_udp_heartbeak_to_all_netcard()
		}
	}
}
