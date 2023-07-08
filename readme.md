### feature

another golang distributed task queue

- ack
- priority
- backoff retry
- at least once
- scheduler
- middleware
- max worker control
- rate limit control
- unique with ttl
- chain
- redis
- sqs

### how to use

#### 1 create a task processor

```go
package test

import (
	"context"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/queue"
	"github.com/pfdtk/goq/task"
	"time"
)

var Conn = "test"

type DemoTask struct {
	task.BaseTask
	logger logger.Logger
}

func NewDemoTask() *DemoTask {
	option := &task.Option{
		Name:      "test",
		OnConnect: Conn,
		QueueType: queue.Redis,
		OnQueue:   "default",
		Status:    task.Active,
		Priority:  0,
		Retries:   0,
		Timeout:   500,
	}
	return &DemoTask{
		BaseTask: task.BaseTask{Option: option},
		logger:   logger.GetLogger(),
	}
}

func (t *DemoTask) Run(_ context.Context, j *task.Job) (any, error) {
	time.Sleep(8 * time.Second)
	t.logger.Info(string(j.RawMessage().Payload))
	return "test", nil
}

func (t *DemoTask) Beforeware() []task.Middleware {
	return []task.Middleware{
		task.NewMaxWorkerLimiter("test", t.GetName(), 1, 10),
		//task.NewLeakyBucketLimiter("test", t.GetName(), redis_rate.PerMinute(1)),
	}
}

func (t *DemoTask) Processware() []task.Middleware {
	return nil
}
```

#### 2 start server and register task

```go
server := NewServer(&ServerConfig{MaxWorker: 5})
// connect
server.AddRedisConnect(test.Conn, test.GetRedis())
server.RegisterTask(test.NewDemoTask())
server.RegisterCronTask("* * * * *", test.NewDemoTask())
server.MustStart(context.Background())
```

#### 3 dispatch message to queue

```go
package goq

import (
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/task"
	"github.com/pfdtk/goq/test"
	"testing"
)

// demo 1
func TestClient(t *testing.T) {
	c := NewClient(&ClientConfig{})
	// connect
	conn := test.GetRedis()
	c.AddRedisConnect(test.Conn, conn)

	err := test.NewDemoTask().Dispatch(
		[]byte("test payload"),
		task.WithDelay(10),
		task.WithPayloadUnique([]byte("test payload"), 10),
	)
	if err != nil {
		t.Error(err)
	}
}

// demo 2
func TestDispatch(t *testing.T) {
	c := NewClient(&ClientConfig{})
	// connect
	conn := test.GetRedis()
	c.AddRedisConnect(test.Conn, conn)
	err := c.Dispatch(test.NewDemoTask(), []byte("test"))
	if err != nil {
		t.Error(err)
	}
}

// demo 3
func TestChain_Dispatch(t *testing.T) {
	// connect
	conn := test.GetRedis()
	connect.AddRedisConnect(test.Conn, conn)
	t1 := test.NewDemoTask().Message([]byte("test chain 1"))
	t2 := test.NewDemoTask().Message([]byte("test chain 2"))
	c := task.NewChain(t1, t2)
	err := c.Dispatch()
	if err != nil {
		t.Error(err)
	}
}
```


### server example

```
./sever_test.go
```

### client example

```
./client_test.go
```
