package main

import (
	"wuwengang/ptunel/handler"
	"wuwengang/ptunel/proxy/http2socks"
	"wuwengang/ptunel/proxy/tcp2socks"
)

func main() {
	go handler.Pipe(http2socks.New())
	handler.Pipe(tcp2socks.New())
}
