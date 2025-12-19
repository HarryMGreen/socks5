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
	ServerIp     string `env:"SERVER_IP" envDefault:""`
	NetInterface string `env:"NET_INTERFACE" envDefault:"eth0"`
}

// GetInterfaceIpv4Addr returns the first IPv4 address found for the specified network interface
func GetInterfaceIpv4Addr(interfaceName string) (string, error) {
	if interfaceName == "" {
		return "", fmt.Errorf("interface name cannot be empty")
	}

	ief, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", fmt.Errorf("interface %s not found: %w", interfaceName, err)
	}

	// Check if interface is up and running
	if ief.Flags&net.FlagUp == 0 {
		return "", fmt.Errorf("interface %s is not up", interfaceName)
	}

	addrs, err := ief.Addrs()
	if err != nil {
		return "", fmt.Errorf("failed to get addresses for interface %s: %w", interfaceName, err)
	}

	// Find first IPv4 address
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipv4 := ipnet.IP.To4(); ipv4 != nil {
				return ipv4.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no IPv4 address found for interface %s", interfaceName)
}

func main() {
	// Working with app params
	cfg := params{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	dockerIp, err := GetInterfaceIpv4Addr(cfg.NetInterface)
	if err != nil {
		log.Printf("[main] get public IP error %s", err.Error())
		dockerIp = "0.0.0.0"
	}
	pubIp := cfg.ServerIp
	if len(pubIp) == 0 {
		log.Printf("[main] server ip is not configured, udp packets may not be received correctly!")
		pubIp = dockerIp
	}

	server, _ := socks5.NewClassicServer(dockerIp+":"+cfg.Port, pubIp, cfg.User, cfg.Password, 60, 60)
	server.ListenAndServe(nil)
}
