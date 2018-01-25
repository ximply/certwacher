package main

import (
	"time"
	"crypto/tls"
	"net"
	"os"
	"fmt"
)

func check(host, port string) (int, error) {
	dialer := &net.Dialer{Timeout: 30 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", host + ":" + port, &tls.Config{
		InsecureSkipVerify: true,
	})
	defer conn.Close()

	if err != nil {
		return -1, err
	}

	if err := conn.Handshake(); err != nil {
		return -1, err
	}

	for _, cert := range conn.ConnectionState().PeerCertificates {
		if cert.IsCA {
			continue
		}

		now := time.Now()

		left := time.Until(cert.NotAfter) - time.Until(now)
		return int(left.Seconds()), nil

	}

	return -1, nil
}


func main() {

	args := os.Args
	if len(os.Args) < 2 {
		fmt.Print("You must take a host argument")
		return
	}

	host := args[1]
	ret, err := check(host, "443")
	if ret < 0 {
		if err != nil {
			fmt.Print("Got a error: " + err.Error())
		} else {
			fmt.Print("Got a error")
		}
	} else {
		fmt.Printf("Remaining time in seconds: %d", ret)
	}
}
