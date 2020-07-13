package store

import (
	"sync"

	"github.com/wndhydrnt/runbooks/pkg/parser"
)

type InMemory struct {
	m    *sync.Mutex
	data map[string]*parser.Runbook
}

func (i *InMemory) Create(rb *parser.Runbook) error {
	i.m.Lock()
	defer i.m.Unlock()
	i.data[rb.Name] = rb
	return nil
}

func (i *InMemory) Delete(name string) error {
	i.m.Lock()
	defer i.m.Unlock()
	delete(i.data, name)
	return nil
}

func (i *InMemory) List() ([]*parser.Runbook, error) {
	result := []*parser.Runbook{}
	i.m.Lock()
	for _, rb := range i.data {
		result = append(result, rb)
	}

	i.m.Unlock()
	return result, nil
}

func NewInMemory() *InMemory {
	return &InMemory{
		m:    &sync.Mutex{},
		data: map[string]*parser.Runbook{},
	}
}
