package main

import (
	lg "crawler/log"
	"crypto/tls"
	"flag"
	"golang.org/x/sys/windows/registry"
	"io"
	"log"
	"net"
	"os/exec"
	"time"
)

var (
	dialer = &net.Dialer{Timeout: 5 * time.Second}
	remote = flag.String("remote", "103.200.6.26:443", "remote tls proxy want tunnel")
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
		InsecureSkipVerify: true,
	}
	remote, err := tls.DialWithDialer(dialer, "tcp", *remote, cfg)
	if err != nil {
		lg.Printf("======== dial remote failed, error %s ========", err)
		return
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

func setupProxy() {
	// set
	key, ok, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	defer key.Close()
	if !ok {
		log.Fatalf("can't create key ProxyServer, %v", err)
	}
	key.SetStringValue("ProxyServer", "127.0.0.1:1443")
	// enable proxy
	key.SetDWordValue("ProxyEnable", uint32(1))
	c := exec.Command("updater.exe")
	err = c.Run()
	if err != nil {
		lg.Printf("Exec error %v", err)
	}
	update()
}
func cleanProxy() {
	key, ok, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	defer key.Close()
	if !ok {
		log.Fatalf("can't cleanProxy,  %v", err)
	}
	// disable proxy
	key.SetDWordValue("ProxyEnable", uint32(0))
	update()
}

func update() {
	c := exec.Command("plugin/updater.exe")
	err := c.Run()
	if err != nil {
		lg.Printf("Exec error %v", err)
	}
}
