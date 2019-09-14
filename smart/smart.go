package smart

const (
	// LOCAL 代表本地网络数据不用转发到代理
	LOCAL = 0x00000001
	// REMOTE 代表远程数据需要走代理
	REMOTE = 0x00000002
)

// Smart 数据分析接口
type Smart interface {
	Parser(func()) //数据解析器
	Local(func())  //本地数据
	Remote(func()) //远程数据
}

// IsLocal 本地数据
func IsLocal(s Smart, p func(l func()), l func()) {
	s.Parser(func() {
		p(l())
	})
	return
}

// IsRemote 远端数据
func IsRemote(s Smart, r func() bool) {
	s.Parser(func() {
		s.Remote(func() {
			isRemote = r()
		})
	})
	return
}
