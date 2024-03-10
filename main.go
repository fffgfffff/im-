package main

import (
	"ginChat/router"
	"ginChat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":9999")
	//addrs, err := net.InterfaceAddrs()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//for _, addr := range addrs {
	//	if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	//		if ipnet.IP.To4() != nil {
	//			fmt.Println("IPv4 Address:", ipnet.IP)
	//		}
	//	}
	//}
}
