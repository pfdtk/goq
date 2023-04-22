package limiter

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"strconv"
)

// Get the Lua script for acquiring a lock.
//
//	KEYS    - The keys that represent available slots
//	ARGV[1] - The limiter name
//	ARGV[2] - The number of seconds the slot should be reserved
//	ARGV[3] - The unique identifier for this lock
var lockScript = redis.NewScript(`
for index, value in pairs(redis.call('mget', unpack(KEYS))) do
    if not value then
        redis.call('set', KEYS[index], ARGV[3], "EX", ARGV[2])
        return ARGV[1]..index
    end
end
`)

// Get the Lua script to atomically release a lock.
//
//	KEYS[1] - The name of the lock
//	ARGV[1] - The unique identifier for this lock
var releaseScript = redis.NewScript(`
if redis.call('get', KEYS[1]) == ARGV[1]
then
    return redis.call('del', KEYS[1])
else
    return 0
end
`)

type Concurrency struct {
	name         string
	redis        *redis.Client
	maxLocks     int
	releaseAfter int
}

func NewConcurrency(
	redis *redis.Client,
	name string,
	maxLocks int,
	releaseAfter int) *Concurrency {

	return &Concurrency{
		redis:        redis,
		name:         name,
		maxLocks:     maxLocks,
		releaseAfter: releaseAfter,
	}
}

func (c *Concurrency) Acquire(id string) (string, error) {
	slots := make([]string, c.maxLocks, c.maxLocks)
	for i := 0; i < c.maxLocks; i++ {
		slots[i] = c.name + strconv.Itoa(i+1)
	}
	keys := slots
	argv := []any{c.name, c.releaseAfter, id}
	val, err := lockScript.Run(context.Background(), c.redis, keys, argv...).Result()
	if err != nil {
		return "", err
	}
	res, err := cast.ToStringE(val)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *Concurrency) Release(key string, id string) bool {
	keys := []string{key}
	argv := []any{id}
	_, err := releaseScript.Run(context.Background(), c.redis, keys, argv...).Result()
	success := err == nil
	return success
}
