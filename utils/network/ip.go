package network

import (
	"errors"
	"net"
	"net/http"
	"strings"

	_ "github.com/Cospk/base-tools/errs"
)

// 定义 HTTP 头部常量。
const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	XClientIP     = "x-client-ip"
)

// GetLocalIP 获取本地IP地址
func GetLocalIP() (string, error) {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	// 遍历每个接口
	var publicIP string
	for _, iface := range interfaces {
		// 检查接口是否启用且不是回环接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取与接口关联的所有地址
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		// 检查每个地址是否为有效的非回环 IPv4 地址
		for _, addr := range addrs {
			// 尝试将地址解析为 IPNet（CIDR 表示法）
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() {
				continue
			}

			ip4 := ipNet.IP.To4()
			if ip4 != nil && !ip4.IsLoopback() {
				// 确保 IP 不是多播地址
				if !ip4.IsMulticast() {
					if !ipNet.IP.IsPrivate() && publicIP == "" {
						// 优先返回内网 IP
						publicIP = ipNet.IP.String()
					} else {
						return ip4.String(), nil
					}

				}
			}
		}
	}

	if publicIP != "" {
		return publicIP, nil
	}
	// 如果没有找到合适的 IP，返回错误
	return "", errors.New("no suitable local IP address found")
}

// GetRpcRegisterIP 获取 RPC 注册 IP
func GetRpcRegisterIP(configIP string) (string, error) {
	registerIP := configIP
	if registerIP == "" {
		ip, err := GetLocalIP()
		if err != nil {
			return "", err
		}
		registerIP = ip
	}
	return registerIP, nil
}

func GetListenIP(configIP string) string {
	if configIP == "" {
		return "0.0.0.0"
	}
	return configIP
}

// RemoteIP returns the remote ip of the request.
func RemoteIP(req *http.Request) string {
	if ip := req.Header.Get(XClientIP); ip != "" {
		return ip
	} else if ip := req.Header.Get(XRealIP); ip != "" {
		return ip
	} else if ip := req.Header.Get(XForwardedFor); ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		ip = req.RemoteAddr
	}

	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}
