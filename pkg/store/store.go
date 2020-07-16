package store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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

type File struct {
	data []*parser.Runbook
}

func (f *File) Create(_ *parser.Runbook) error {
	return fmt.Errorf("store File does not support Create operation")
}

func (f *File) Delete(_ string) error {
	return fmt.Errorf("store File does not support Delete operation")
}

func (f *File) List() ([]*parser.Runbook, error) {
	return f.data, nil
}

func NewFile(pathPattern string, rbp *parser.Parser) (*File, error) {
	paths, err := filepath.Glob(pathPattern)
	if err != nil {
		return nil, err
	}

	f := &File{}
	for _, p := range paths {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}

		rb, err := rbp.ParseRunbook(b)
		if err != nil {
			return nil, err
		}

		f.data = append(f.data, &rb)
	}

	return f, nil
}
