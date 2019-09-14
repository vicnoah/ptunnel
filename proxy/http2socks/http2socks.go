package http2socks

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	httparser "wuwengang/ptunel/parser/http"
	"wuwengang/ptunel/smart"

	"golang.org/x/net/proxy"
)

// New 新建代理链
func New() *HTTP2Socks {
	return &HTTP2Socks{}
}

// HTTP2Socks http <->   transfer   <-> socks 代理链
type HTTP2Socks struct{}

// Listen 帧听端口
func (h *HTTP2Socks) Listen() {
	addr := ":8888"
	// listen http
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.In(func() {
				h.Out(func() {
					outFunc(w, r, h)
				})
			})
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Println(fmt.Sprintf("http2socks->bind: %s, start listening...\r\n", addr))
	server.ListenAndServe()
}

// In 数据入站
func (h *HTTP2Socks) In(inFunc func()) {
	inFunc()
}

// Transport 数据传输者
func (h *HTTP2Socks) Transport(trans func()) {
	trans()
}

// Out 数据出站
func (h *HTTP2Socks) Out(outFunc func()) {
	// outer
	outFunc()
}

func outFunc(w http.ResponseWriter, r *http.Request, h *HTTP2Socks) {
	// 数据解析
	smart.IsLocal(httparser.New(httparser.P()))
	// 数据分类型代理
	if r.Method == http.MethodConnect {
		// Connect连接
		h.Transport(func() {
			handleSOCKSTunnel(w, r)
		})
	} else {
		h.Transport(func() {
			handleHTTProxy(w, r)
		})
	}
}

// Start 启动代理
func (h *HTTP2Socks) Start() {
	// In -> Transfer
	// Transfer -> In
	h.Listen()
}

func handleHTTProxy(w http.ResponseWriter, r *http.Request) {
	dialer, err := proxy.SOCKS5("tcp", "192.168.5.50:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Printf("can't connect to the proxy:%v\r\n", err)
		return
	}
	tr := &http.Transport{Dial: dialer.Dial}
	httpClient := &http.Client{Transport: tr}

	req, _ := http.NewRequest(r.Method, r.RequestURI, r.Body)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func handleSOCKSTunnel(w http.ResponseWriter, r *http.Request) {
	dialer, err := proxy.SOCKS5("tcp", "192.168.5.50:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Printf("can't connect to the proxy:%v\r\n", err)
	}
	// ssocks5Proxy := socks.DialSocksProxy(socks.SOCKS5, "192.168.5.50:1080")
	destConn, err := dialer.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
