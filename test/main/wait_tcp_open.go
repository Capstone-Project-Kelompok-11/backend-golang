package main

import (
	"fmt"
	"net"
	"time"
)

func WaitTcpOpened(address string, timeout time.Duration) bool {

	Interval := time.Second * 1
	ReachTimeOut := time.Second * 0

	for ReachTimeOut < timeout {

		var err error
		var conn net.Conn
		var opened bool

		opened = true

		if conn, err = net.DialTimeout("tcp", address, timeout); err != nil {

			fmt.Println("Connecting error:", err)
			opened = false
		}

		if conn != nil {

			defer conn.Close()
			if opened {

				fmt.Println("Opened", address)
				return true
			}
		}

		ReachTimeOut += Interval
		time.Sleep(Interval)
	}

	fmt.Println("Timeout has been reached")
	return false
}

func main() {

	WaitTcpOpened(net.JoinHostPort("localhost", "8000"), time.Minute*1)
}
