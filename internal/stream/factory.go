package stream

import (
	"os"
	"strconv"
	"time"

	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// Stream manager types
const (
	TypeMemory = "memory"
	TypeRedis  = "redis"
)

// NewStreamManager creates a stream manager
func NewStreamManager() (interfaces.StreamManager, error) {
	switch os.Getenv("STREAM_MANAGER_TYPE") {
	case TypeRedis:
		db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			db = 0
		}
		ttl := time.Hour // Default 1 hour
		return NewRedisStreamManager(
			os.Getenv("REDIS_ADDR"),
			os.Getenv("REDIS_PASSWORD"),
			db,
			os.Getenv("REDIS_PREFIX"),
			ttl,
		)
	default:
		return NewMemoryStreamManager(), nil
	}
}
