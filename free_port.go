package Gorra

import (
	"fmt"
	"net"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = l.Close()
	}()

	fmt.Println("[Gorra] Get Free Port Success")

	return l.Addr().(*net.TCPAddr).Port, nil
}
