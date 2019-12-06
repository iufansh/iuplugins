package queue

import (
	"fmt"
)

type Queue interface {
	// get value by default key.
	Get() interface{}
	// put into queue
	Put(val interface{}) error
	// get value by name
	GetByName(name string) interface{}
	// put value by name
	PutByName(name string, val interface{}) error
	// clear all Queue.
	ClearAll() error
	// start gc routine based on config string settings.
	StartAndGC(config string) error
}

// Instance is a function create a new Cache Instance
type Instance func() Queue

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
// support list and memory
func NewQueue(adapterName, config string) (adapter Queue, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("queue: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
