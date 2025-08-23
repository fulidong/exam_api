package isnowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	// 位数分配定义
	workerBits   = 10 // 节点ID位数
	sequenceBits = 12 // 每毫秒内ID生成位数

	// 最大值掩码
	maxWorkerID = -1 ^ (-1 << workerBits)   // 用于防止越界
	maxSequence = -1 ^ (-1 << sequenceBits) // 同上

	// 时间偏移量（毫秒）
	workerShift    = sequenceBits
	timestampShift = workerBits + sequenceBits
	sequenceMask   = maxSequence

	// 毫秒内最大生成数量
	sequenceMax = 4095 // maxSequence = 4095
)

var SnowFlake *Snowflake
var (
	ErrTimeBackwards     = errors.New("clock moved backwards")
	ErrSequenceExhausted = errors.New("sequence overflow in the same millisecond")
	ErrWorkerIDInvalid   = errors.New("worker ID is invalid")
)

type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
}

func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrWorkerIDInvalid
	}
	return &Snowflake{
		workerID:  workerID,
		timestamp: time.Now().UnixNano() / 1e6,
	}, nil
}

func (sf *Snowflake) NextID(key string) (string, error) {
	id, err := sf.nextID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%d", key, id), nil
}

func (sf *Snowflake) nextID() (int64, error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	now := time.Now().UnixNano() / 1e6 // 当前毫秒时间戳

	if now < sf.timestamp {
		return 0, ErrTimeBackwards
	}

	if now == sf.timestamp {
		sf.sequence = (sf.sequence + 1) & sequenceMask
		if sf.sequence == 0 {
			// 同一毫秒内序号已满，等待下一毫秒
			<-time.After(time.Millisecond)
			return sf.nextID()
		}
	} else {
		sf.sequence = 0
	}

	sf.timestamp = now

	id := (now << timestampShift) |
		(sf.workerID << workerShift) |
		sf.sequence

	return id, nil
}
