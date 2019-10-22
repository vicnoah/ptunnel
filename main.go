package main

import (
	"github.com/wuwengang/ptunnel/handler"
	"github.com/wuwengang/ptunnel/proxy/http2socks"
	"github.com/wuwengang/ptunnel/proxy/tcp2socks"
)

func main() {
	go handler.Pipe(http2socks.New())
	handler.Pipe(tcp2socks.New())
}
