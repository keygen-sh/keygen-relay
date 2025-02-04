package locker

import (
	"fmt"
	"net"
)

func getPrivateIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to read network interfaces: %w", err)
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 { // skip down and loopback interfaces
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			return "", fmt.Errorf("failed to read network interface: %w", err)
		}

		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok && ip.IP.To4() != nil {
				if ip.IP.IsPrivate() {
					return ip.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("no private ip found")
}
