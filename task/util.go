package task

import (
	"github.com/pfdtk/goq/queue"
	"sort"
	"sync"
)

func SortTask(tasks *sync.Map) []Task {
	var pairs []Task
	tasks.Range(func(key, value any) bool {
		pairs = append(pairs, value.(Task))
		return true
	})
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Priority() < pairs[j].Priority()
	})
	return pairs
}

func GetRedisTask(tasks *sync.Map) []Task {
	var pairs []Task
	tasks.Range(func(key, value any) bool {
		v := value.(Task)
		if v.QueueType() != queue.Redis {
			return true
		}
		pairs = append(pairs, v)
		return true
	})
	return pairs
}
