package main

import (
	"github.com/vicnoah/ptunnel/handler"
	"github.com/vicnoah/ptunnel/proxy/http2socks"
	"github.com/vicnoah/ptunnel/proxy/tcp2socks"
)

func main() {
	go handler.Pipe(http2socks.New())
	handler.Pipe(tcp2socks.New())
}
