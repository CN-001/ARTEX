package db

import "sync"

// ResetMasterKeyCacheForTest 重置主密钥缓存（loadMasterKey 用 sync.Once 缓存）。
//
// 之所以放在非 _test.go 文件里：其它包的测试（如 server 的归档加密测试）也需要在
// 「有主密钥 / 无主密钥」两种状态间切换，而 _test.go 里的导出符号对其他包不可见。
// 生产代码不应调用它。
func ResetMasterKeyCacheForTest() {
	masterOnce = sync.Once{}
	masterKey = nil
	masterErr = nil
}