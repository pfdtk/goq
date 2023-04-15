package utils

import (
	"github.com/pfdtk/goq/base"
	"github.com/pfdtk/goq/task"
	"sort"
	"sync"
)

func SortTask(tasks *sync.Map) []task.Task {
	var pairs []task.Task
	tasks.Range(func(key, value any) bool {
		pairs = append(pairs, value.(task.Task))
		return true
	})
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Priority() < pairs[j].Priority()
	})
	return pairs
}

func GetRedisTask(tasks *sync.Map) []task.Task {
	var pairs []task.Task
	tasks.Range(func(key, value any) bool {
		v := value.(task.Task)
		if v.QueueType() != base.Redis {
			return true
		}
		pairs = append(pairs, v)
		return true
	})
	return pairs
}
