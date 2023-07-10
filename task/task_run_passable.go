package task

type RunPassable struct {
	task Task
	job  *Job
}

func NewRunPassable(t Task, j *Job) *RunPassable {
	return &RunPassable{task: t, job: j}
}
