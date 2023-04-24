package task

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
