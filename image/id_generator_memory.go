package image

import (
	"fmt"
	"sync"
)

type IdGeneratorMemory struct {
	l sync.Mutex

	nextId int64
}

// TODO: consider generating hash
// to make service stateless
func (i *IdGeneratorMemory) Generate() (string, error) {
	i.l.Lock()
	defer i.l.Unlock()

	returnId := i.nextId
	i.nextId++

	return fmt.Sprintf("%d", returnId), nil
}
