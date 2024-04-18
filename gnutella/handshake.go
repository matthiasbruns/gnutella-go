package gnutella

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	tcpHandshakeTimeout = 2 * time.Second
	tcpWriteTimeout     = 1 * time.Second
)

func InitTCPHandshake(ctx context.Context, listen, target *net.TCPAddr, success chan<- bool) (*net.Conn, error) {
	fmt.Println("Initiating TCP handshake with", target.String())

	conn, err := net.DialTimeout("tcp", target.String(), tcpHandshakeTimeout)
	if err != nil {
		fmt.Println("Error dialing TCP:", err)
		return nil, err
	}

	// stage 1
	payload := fmt.Sprintf("GNUTELLA CONNECT/0.6\nListen-IP: %s\nRemote-IP: %s\nUser-Agent: gnutella-go 0.0.1\nAccept: application/x-gnutella2\nX-Hub: False", listen, target.IP)
	fmt.Println("Sending stage 1:", payload)

	go func() {
		data := []byte(payload)
		n, err := conn.Write(data)
		if err != nil {
			fmt.Println("Error while sending handshake stage 1", err)
			success <- false
			return
		}

		if n != len(data) {
			fmt.Println("Error while sending handshake stage 1: not all bytes sent")
			success <- false
			return
		}

		fmt.Println("Sent stage 1 successfully")
		success <- true
	}()

	return nil, nil
}
