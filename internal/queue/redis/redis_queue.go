package redis

import (
	"context"
	"encoding/json"
	"github.com/pfdtk/goq/common"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"time"
)

// Get the Lua script for popping the next job off of the queue.
//
//	KEYS[1] - The queue to pop jobs from, for example: queues:foo
//	KEYS[2] - The queue to place reserved jobs on, for example: queues:foo:reserved
//	ARGV[1] - The current UNIX timestamp
var popScript = redis.NewScript(`
-- Pop the first job off of the queue...
local job = redis.call('lpop', KEYS[1])
local reserved = false

if(job ~= false) then
    -- Increment the attempt count and place job on the reserved queue...
    reserved = cjson.decode(job)
		local timeout = reserved['timeout']
		local visibility_timeout = tonumber(ARGV[1]) + tonumber(timeout)
    reserved['attempts'] = reserved['attempts'] + 1
    reserved = cjson.encode(reserved)
    redis.call('zadd', KEYS[2], visibility_timeout, reserved)
end

return {job, reserved}
`)

// Get the Lua script to migrate expired jobs back onto the queue.
//
//	KEYS[1] - The queue we are removing jobs from, for example: queues:foo:reserved
//	KEYS[2] - The queue we are moving jobs to, for example: queues:foo
//	ARGV[1] - The current UNIX timestamp
var migrateScript = redis.NewScript(`
-- Get all of the jobs with an expired "score"...
local val = redis.call('zrangebyscore', KEYS[1], '-inf', ARGV[1])

-- If we have values in the array, we will remove them from the first queue
-- and add them onto the destination queue in chunks of 100, which moves
-- all of the appropriate jobs onto the destination queue very safely.
if(next(val) ~= nil) then
    redis.call('zremrangebyrank', KEYS[1], 0, #val - 1)

    for i = 1, #val, 100 do
        redis.call('rpush', KEYS[2], unpack(val, i, math.min(i+99, #val)))
    end
end

return val
`)

type Queue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *Queue {
	return &Queue{client: client}
}

func (r Queue) Size(ctx context.Context, queue string) (int64, error) {
	size, err := r.client.LLen(ctx, queue).Result()
	return size, err
}

func (r Queue) Push(ctx context.Context, message *common.Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = r.client.RPush(ctx, message.Queue, bytes).Result()
	return err
}

func (r Queue) Later(ctx context.Context, message *common.Message, at time.Time) error {
	queue := r.GetDelayedKey(message.Queue)
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	score := at.Unix()
	_, err = r.client.ZAdd(ctx, queue, redis.Z{
		Score:  float64(score),
		Member: bytes,
	}).Result()

	return err
}

func (r Queue) Pop(ctx context.Context, queue string) (*common.Message, error) {
	keys := []string{queue, r.GetReservedKey(queue)}
	argv := []any{time.Now().Unix()}
	val, err := popScript.Run(ctx, r.client, keys, argv...).Result()
	if err != nil {
		return nil, err
	}
	res, err := cast.ToStringSliceE(val)
	if err != nil {
		return nil, err
	}
	msg := common.Message{}
	err = json.Unmarshal([]byte(res[0]), &msg)
	if err != nil {
		return nil, err
	}
	msg.Reserved = res[1]
	return &msg, nil
}

func (r Queue) Migrate(ctx context.Context, from string, to string) error {
	keys := []string{from, to}
	argv := []any{time.Now().Unix()}
	_, err := migrateScript.Run(ctx, r.client, keys, argv...).Result()
	return err
}

func (r Queue) GetReservedKey(queue string) string {
	return queue + ":reserved"
}

func (r Queue) GetDelayedKey(queue string) string {
	return queue + ":delayed"
}
