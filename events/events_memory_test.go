package events_test

import (
	"github.com/progimage/events"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEmitEvent(t *testing.T) {
	id := "1"

	memoryEvents := events.InMemoryEvents{
		Store:    map[string][]events.State{},
		Watchers: map[string][]chan events.State{},
	}
	state := events.State{
		Val:       "some-val",
		Timestamp: time.Time{},
	}
	memoryEvents.Emit(id, state)
	rState, found := memoryEvents.Get(id)

	require.True(t, found)
	require.Len(t, rState, 1)
	require.Equal(t, rState[0], state)

	t.Run("returns false if id not found", func(t *testing.T) {
		_, found := memoryEvents.Get("non-existent")
		require.False(t, found)
	})

	t.Run("watch for events", func(t *testing.T) {
		watch := memoryEvents.Watch(id)
		require.True(t, found)
		memoryEvents.Emit(id, state)

		var watchedState events.State
		require.Eventually(t, func() bool {
			watchedState = <-watch
			return true
		}, 3*time.Second, 1*time.Second)

		require.Equal(t, watchedState, state)
	})
}
