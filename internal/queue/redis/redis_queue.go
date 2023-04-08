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
//	ARGV[1] - Current timestamp
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
	queue := message.Queue + ":delayed"
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
	keys := []string{queue, queue + ":reserved"}
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
