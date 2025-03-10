package argon2

import (
	"testing"
)

// 性能测试函数
func BenchmarkArgon2Hash(b *testing.B) {
	password := "execmt111"
	salt := "exec"
	// 重置计时器，避免生成盐值的时间影响测试结果
	b.ResetTimer()

	// 运行 b.N 次，测试生成哈希值的速度
	for i := 0; i < b.N; i++ {
		GetEncryptString(password, salt)
	}
}
