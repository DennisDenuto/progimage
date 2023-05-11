package image_test

import (
	"github.com/progimage/image"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestReturnsIncreasingID(t *testing.T) {
	idGeneratorMemory := image.IdGeneratorMemory{}

	id, err := idGeneratorMemory.Generate()
	require.NoError(t, err)
	require.Equal(t, id, "0")

	id, err = idGeneratorMemory.Generate()
	require.NoError(t, err)
	require.Equal(t, id, "1")
}

func TestGeneratingIDInParallel(t *testing.T) {
	idGeneratorMemory := image.IdGeneratorMemory{}

	N := 10
	generatedIds := make(chan string, N)
	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := idGeneratorMemory.Generate()
			require.NoError(t, err)
			generatedIds <- id
		}()
	}
	require.Eventually(t, func() bool {
		wg.Wait()
		return true
	}, 5*time.Second, 1*time.Second)
	close(generatedIds)
	assertUnique(t, generatedIds)
}

func assertUnique(t *testing.T, ids chan string) {
	t.Helper()
	visited := map[string]string{}

	for v := range ids {
		if _, found := visited[v]; found {
			require.Fail(t, "found duplicate ID")
		}
		visited[v] = v
	}
}
