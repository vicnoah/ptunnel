package tcp2socks

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/net/proxy"
)

const (
	_PROXY_PROTOCOL = "tcp"
)

// New 新建代理链
func New() *TCP2Socks {
	return &TCP2Socks{}
}

// TCP2Socks tcp <->   transfer   <-> socks 代理链
type TCP2Socks struct{}

// Listen 监听地址
func (h *TCP2Socks) Listen() {
	// listen tcp
	// 绑定监听地址
	addr := ":8889"
	listener, err := net.Listen(_PROXY_PROTOCOL, addr)
	if err != nil {
		log.Fatal(fmt.Sprintf("listen err: %v", err))
	}
	defer listener.Close()
	log.Println(fmt.Sprintf("tcp2socks->bind: %s, start listening...\r\n", addr))

	for {
		// Accept 会一直阻塞直到有新的连接建立或者listen中断才会返回
		conn, err := listener.Accept()
		if err != nil {
			// 通常是由于listener被关闭无法继续监听导致的错误
			log.Fatal(fmt.Sprintf("accept err: %v", err))
		}
		// 开启新的 goroutine 处理该连接
		go h.In(func() {
			h.Out(func() {
				outFunc(conn, h)
			})
		})
	}
}

// In 数据入站
func (h *TCP2Socks) In(inFunc func()) {
	// tcp
	inFunc()
}

// Out 数据出站
func (h *TCP2Socks) Out(outFunc func()) {
	outFunc()
}

func outFunc(conn net.Conn, h *TCP2Socks) {
	defer conn.Close()
	// socks5
	dialer, err := proxy.SOCKS5("tcp", "192.168.5.50:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Printf("can't connect to the proxy:%v\r\n", err)
		conn.Close()
		return
	}
	dstConn, err := dialer.Dial("tcp", conn.RemoteAddr().String())
	if err != nil {
		fmt.Println("error:", err)
		dstConn.Close()
		return
	}
	defer dstConn.Close()
	h.Transport(func() {
		transport(dstConn, conn)
	})
}

// Transport 数据传输方法
func (h *TCP2Socks) Transport(trans func()) {
	trans()
}

func transport(dst net.Conn, src net.Conn) {
	defer func() {
		fmt.Println("关闭连接")
	}()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(src, dst)
		wg.Done()
	}()

	go func() {
		io.Copy(dst, src)
		wg.Done()
	}()

	wg.Wait()
	return
}

// Start 启动代理
func (h *TCP2Socks) Start() {
	// In -> Transfer
	// Transfer -> In
	h.Listen()
}
