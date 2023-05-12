package events

import "sync"

type InMemoryEvents struct {
	l              sync.Mutex
	Store          map[string][]State
	Watchers       map[string][]chan State
	GlobalWatchers []chan State
}

func NewInMemoryEvents() *InMemoryEvents {
	return &InMemoryEvents{
		Store:          map[string][]State{},
		Watchers:       map[string][]chan State{},
		GlobalWatchers: []chan State{},
	}
}
func (i *InMemoryEvents) Get(id string) ([]State, bool) {
	v, found := i.Store[id]
	return v, found
}

func (i *InMemoryEvents) Emit(id string, state State) {
	i.l.Lock()
	defer i.l.Unlock()

	i.Store[id] = append(i.Store[id], state)
	i.broadcast(id, state)
}

func (i *InMemoryEvents) Watch(id string) <-chan State {
	i.l.Lock()
	defer i.l.Unlock()

	ch := make(chan State)
	i.Watchers[id] = append(i.Watchers[id], ch)
	return ch
}

func (i *InMemoryEvents) WatchAll() <-chan State {
	i.l.Lock()
	defer i.l.Unlock()

	ch := make(chan State)
	i.GlobalWatchers = append(i.GlobalWatchers, ch)
	return ch
}

func (i *InMemoryEvents) broadcast(id string, state State) {
	for _, w := range i.Watchers[id] {
		go func(w chan State) {
			w <- state
		}(w)
	}

	for _, w := range i.GlobalWatchers {
		go func(w chan State) {
			w <- state
		}(w)
	}
}
