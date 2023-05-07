### feature

another golang distributed queue

- job ack
- queue priority
- job backoff retry
- job execute timeout notice by event dispatch
- job recover by ack feature
- task scheduler
- task run through middleware
- max worker of each task control by middleware
- rate limit control by middleware
- job unique with ttl
- redis
- sqs [todo]

### server example

```
./sever_test.go
```

### client example

```
./client_test.go
```
