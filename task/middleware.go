package task

import "github.com/pfdtk/goq/pipeline"

type Middleware pipeline.Handler

func CastMiddleware(m []Middleware) []pipeline.Handler {
	h := make([]pipeline.Handler, len(m), len(m))
	for i, v := range m {
		h[i] = pipeline.Handler(v)
	}
	return h
}
