package goq

type MigrateErrorEvent struct {
	err error
}

func NewMigrateErrorEvent(e error) *MigrateErrorEvent {
	return &MigrateErrorEvent{err: e}
}

func (j *MigrateErrorEvent) Name() string {
	return "MigrateErrorEvent"
}

func (j *MigrateErrorEvent) Value() any {
	return j.err
}
