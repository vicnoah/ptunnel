// Package handler 用于处理代理链
package handler

// Tunel 代理隧道
// 数据流 In
//          -> Transfer
//                      -> Out
//             Transfer <-
//       In <-
/* type Tunel interface {
	In()               //隧道数据入站方法
	Out()              //隧道数据出站方法
	Transfer(f func()) //数据传输者(中间人)
	Start()            //打开Tunel
} */
type Tunel interface {
	Listen()                //帧听方法
	In(inFunc func())       //隧道数据入站方法
	Out(outFunc func())     //隧道数据出站方法
	Transport(trans func()) //数据传输方法
	Start()                 //打开Tunel
}

// Pipe 启动代理管道
func Pipe(t Tunel) {
	t.Start()
}
