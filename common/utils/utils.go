package utils

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

func ObjName(t any) string {
	typeOf := reflect.TypeOf(t)
	kind := typeOf.Kind()
	name := typeOf.String()
	if kind == reflect.Pointer {
		name = typeOf.Elem().String()
	}
	return name
}

func GetIps() []string {
	addrList, err := net.InterfaceAddrs()
	ips := make([]string, 0)
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ips
	}
	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func GetRealIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(addr.String(), ":")[0]
	return ip
}
