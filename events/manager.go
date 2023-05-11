package events

import "time"

type State struct {
	Val       string
	Timestamp time.Time
	Metadata  map[string]string
}

type FileEvents interface {
	Emit(id string, state State)
	Watch(id string) <-chan State
}
