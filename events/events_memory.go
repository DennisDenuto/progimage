package events

type InMemoryEvents struct {
	Store map[string][]State
}

func (i *InMemoryEvents) Emit(id string, state State) {
	i.Store[id] = append(i.Store[id], state)
}

func (i *InMemoryEvents) Watch(id string) <-chan State {
	//TODO implement me
	panic("implement me")
}
