package task

type JobExecTimeoutEventValue struct {
	job *Job
}

type JobExecTimeoutEvent struct {
	value *JobExecTimeoutEventValue
}

func NewJobExecTimeoutEvent(job *Job) *JobExecTimeoutEvent {
	return &JobExecTimeoutEvent{
		&JobExecTimeoutEventValue{job: job},
	}
}

func (j *JobExecTimeoutEvent) Name() string {
	return "JobExecTimeoutEvent"
}

func (j *JobExecTimeoutEvent) Value() any {
	return j.value
}
