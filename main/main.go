package main

import (
	"fmt"
	"log"
	"net"

	"github.com/caarlos0/env/v6"
	"github.com/txthinking/socks5"
)

type params struct {
	User         string `env:"PROXY_USER" envDefault:""`
	Password     string `env:"PROXY_PASSWORD" envDefault:""`
	Port         string `env:"PROXY_PORT" envDefault:"1080"`
	NetInterface string `env:"NET_INTERFACE" envDefault:"eth0"`
}

func GetInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return "", fmt.Errorf("interface %s don't have an ipv4 address", interfaceName)
	}
	return ipv4Addr.String(), nil
}

func main() {
	// Working with app params
	cfg := params{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	pubIp, err := GetInterfaceIpv4Addr(cfg.NetInterface)
	if err != nil {
		log.Printf("get publice IP error %s", err.Error())
		pubIp = "0.0.0.0"
	}

	server, _ := socks5.NewClassicServer(pubIp+":"+cfg.Port, pubIp, cfg.User, cfg.Password, 60, 60)
	server.ListenAndServe(nil)
}
