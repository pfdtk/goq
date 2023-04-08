package utils

import (
	"github.com/pfdtk/goq/iface"
	"sort"
	"sync"
)

func SortTask(tasks *sync.Map) []iface.Task {
	var pairs []iface.Task
	tasks.Range(func(key, value any) bool {
		pairs = append(pairs, value.(iface.Task))
		return true
	})
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Priority() < pairs[j].Priority()
	})
	return pairs
}
