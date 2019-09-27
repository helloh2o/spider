package main

import (
	lg "crawler/log"
	"crypto/tls"
	"flag"
	"io"
	"net"
	"time"
)

var (
	dialer = &net.Dialer{Timeout: 5 * time.Second}
	remote = flag.String("remote", "", "remote tls proxy want tunnel")
	host   string
)

func main() {
	flag.Parse()
	if *remote == "" {
		panic("input the remote tls proxy want tunnel, eg: xx0proxy.info:443")
	}
	var err error
	host, _, err = net.SplitHostPort(*remote)
	if err != nil {
		panic(err)
	}
	local := ":1443"
	ln, err := net.Listen("tcp", local)
	if err != nil {
		panic(err)
	}
	lg.Printf("http tunnel listen on %s", local)
	for {
		con, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleRequest(con)
	}
}

func handleRequest(src net.Conn) {
	cfg := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}
	remote, err := tls.DialWithDialer(dialer, "tcp", cfg.ServerName+":443", cfg)
	if err != nil {
		lg.Printf("======== dial remote failed, error %s ========", err)
	}
	tunnel(src, remote)
}
func tunnel(src, remote net.Conn) {
	defer func() {
		src.Close()
		remote.Close()
	}()
	done := make(chan struct{})
	go func() {
		io.Copy(remote, src)
		done <- struct{}{}

	}()
	go func() {
		io.Copy(src, remote)
		done <- struct{}{}
	}()
	<-done
}
