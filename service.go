package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
)

const (
	cert = "./certs/steamcommunity.crt"
	key  = "./certs/steamcommunity.key"
)

func startService(ip string, logger *log.Logger) {
	cfg := &tls.Config{}
	cfg.Certificates = make([]tls.Certificate, 1)
	var err error
	cfg.Certificates[0], err = tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logger.Printf("load cert error is %s", err.Error())
		return
	}
	ls, err := tls.Listen("tcp", ":443", cfg)
	if err != nil {
		logger.Printf("listen local error %s", err.Error())
		return
	}
	for {
		conn, err := ls.Accept()
		if err != nil {
			logger.Printf("accept connection error:%s", err.Error())
			break
		}

		go handleConn(conn, ip, logger)
	}
}

func handleConn(conn net.Conn, ip string, logger *log.Logger) {
	defer conn.Close()

	cfg := &tls.Config{
		ServerName:         "",
		InsecureSkipVerify: true,
	}

	remote, err := tls.Dial("tcp", ip+":443", cfg)
	if err != nil {
		logger.Printf("dial remote error %s", err.Error())
		return
	}
	defer remote.Close()

	go io.Copy(remote, conn)
	io.Copy(conn, remote)
}
