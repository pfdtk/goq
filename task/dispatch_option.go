package task

import (
	"crypto/md5"
	"fmt"
)

type DispatchOpt struct {
	Delay     int64
	UniqueId  string
	UniqueTTL int64
}

type DispatchOptFunc func(opt *DispatchOpt)

func WithDelay(delay int64) DispatchOptFunc {
	return func(opt *DispatchOpt) {
		opt.Delay = delay
	}
}

func WithUnique(uniqueId string, uniqueTTL int64) DispatchOptFunc {
	return func(opt *DispatchOpt) {
		opt.UniqueId = uniqueId
		opt.UniqueTTL = uniqueTTL
	}
}

func WithPayloadUnique(payload []byte, uniqueTTL int64) DispatchOptFunc {
	sum := md5.Sum(payload)
	uniqueId := fmt.Sprintf("%x", sum)
	return func(opt *DispatchOpt) {
		opt.UniqueId = uniqueId
		opt.UniqueTTL = uniqueTTL
	}
}
