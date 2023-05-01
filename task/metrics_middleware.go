package task

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	jobProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "job_processed",
			Help: "The total number of processed job",
		},
		[]string{"task"},
	)

	jobFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "job_failed",
			Help: "The total number of times processing failed",
		},
		[]string{"task"},
	)

	jobInProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "job_in_progress",
			Help: "The number of jobs currently being processed",
		},
		[]string{"task"},
	)
)

func NewMetricsMiddleware() Middleware {
	return func(p any, next func(p any) any) any {
		pp, ok := p.(*RunPassable)
		if !ok {
			return next(p)
		}
		jobInProgress.WithLabelValues(pp.task.GetName()).Inc()
		r := next(p)
		jobInProgress.WithLabelValues(pp.task.GetName()).Dec()

		err, ok := p.(error)
		if ok && err != nil {
			jobFailed.WithLabelValues(pp.task.GetName()).Inc()
		}

		jobProcessed.WithLabelValues(pp.task.GetName()).Inc()

		return r
	}
}
