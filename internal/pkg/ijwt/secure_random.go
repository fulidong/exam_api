package ijwt

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"time"
)

// SecureInt63 获取密码学安全的随机int63
func SecureInt63() int64 {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err == nil {
		// 转换为int64并确保是正数
		return int64(binary.BigEndian.Uint64(buf[:]) & math.MaxInt64)
	}

	// 简单回退方案：使用纳秒时间戳
	return time.Now().UnixNano()
}

// SecureBytes 获取密码学安全的随机字节
func SecureBytes(n int) []byte {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err == nil {
		return buf
	}

	// 简单回退方案：使用时间戳生成伪随机字节
	ts := time.Now().UnixNano()
	for i := range buf {
		// 简单的伪随机算法
		ts = (ts*1103515245 + 12345) & 0x7FFFFFFFFFFFFFFF
		buf[i] = byte(ts % 256)
	}
	return buf
}
