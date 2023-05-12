package events

import "time"

type State struct {
	ID        string
	Val       string
	Timestamp time.Time
	Metadata  map[string]string
}

type FileEvents interface {
	Emit(id string, state State)
	Get(id string) ([]State, bool)
	Watch(id string) <-chan State
	WatchAll() <-chan State
}
