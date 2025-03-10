package argon2

import (
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

const (
	time    uint32 = 2        //迭代次数
	memory  uint32 = 8 * 1024 //内存使用量（8MB）
	threads uint8  = 4        //并行度
	keyLen  uint32 = 32       //生成哈希值长度（32字节）
)

func GetEncryptString(value, salt string) string {
	//生成哈希值
	hashed := argon2.IDKey([]byte(value), []byte(salt), time, memory, threads, keyLen)
	//哈希值转化为字符串
	return base64.RawStdEncoding.EncodeToString(hashed)
}
