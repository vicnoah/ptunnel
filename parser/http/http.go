package httparser

import "net/http"

// New 初始化http解析器
func New() *HTTParser {
	return &HTTParser{}
}

// HTTParser http解析器
type HTTParser struct{}

// Parser 解析
func (h *HTTParser) Parser(p func()) {
	p()
}

// Local 本地地址
func (h *HTTParser) Local(l func()) {
	l()
}

// Remote 远程地址
func (h *HTTParser) Remote(r func()) {
	r()
}

// P 解析
func P(w http.ResponseWriter, r *http.Request) {

}

// L 本地地址判断
func L() bool {
	return false
}

// R 远程地址判断
func R() bool {
	return false
}
