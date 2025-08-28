package _const

var (
	GetQuestionsBySalesPaperIdRedisKey = "questions:%s"        // %s为试卷id
	RedisLockKey                       = "exam_lock:submit:%s" // 分布式锁 key
	RedisSubmitKey                     = "exam_submitted:%s"   // 已提交标记 key
	UnlockScript                       = `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
		`
)
