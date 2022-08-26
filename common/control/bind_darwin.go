package control

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func NewBindManager() BindManager {
	return &myBindManager{
		interfaceIndexByName: make(map[string]int),
	}
}

func BindToInterface(manager BindManager, interfaceName string) Func {
	return func(network, address string, conn syscall.RawConn) error {
		index, err := manager.IndexByName(interfaceName)
		if err != nil {
			return err
		}
		return bindToInterface(conn, network, index)
	}
}

func BindToInterfaceFunc(manager BindManager, interfaceNameFunc func(network, address string) string) Func {
	return func(network, address string, conn syscall.RawConn) error {
		interfaceName := interfaceNameFunc(network, address)
		if interfaceName == "" {
			return nil
		}
		index, err := manager.IndexByName(interfaceName)
		if err != nil {
			return err
		}
		return bindToInterface(conn, network, index)
	}
}

func BindToInterfaceIndexFunc(interfaceIndexFunc func(network, address string) int) Func {
	return func(network, address string, conn syscall.RawConn) error {
		index := interfaceIndexFunc(network, address)
		return bindToInterface(conn, network, index)
	}
}

func bindToInterface(conn syscall.RawConn, network string, index int) error {
	return Raw(conn, func(fd uintptr) error {
		switch network {
		case "tcp6", "udp6":
			return unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_BOUND_IF, index)
		default:
			return unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_BOUND_IF, index)
		}
	})
}
