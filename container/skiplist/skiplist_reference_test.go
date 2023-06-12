// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package skiplist

import (
	"sync"
	"sync/atomic"
)

// This file contains reference map implementations for unit-tests.

// mapInterface is the interface Map implements.
type mapInterface interface {
	Find(int) (int, bool)
	Insert(key, value int)
	Delete(int)
}

var (
	_ mapInterface = &RWMutexMap{}
	_ mapInterface = &DeepCopyMap{}
)

// RWMutexMap is an implementation of mapInterface using a sync.RWMutex.
type RWMutexMap struct {
	mu    sync.RWMutex
	dirty map[int]int
}

func (m *RWMutexMap) Find(key int) (value int, ok bool) {
	m.mu.RLock()
	value, ok = m.dirty[key]
	m.mu.RUnlock()
	return
}

func (m *RWMutexMap) Insert(key, value int) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[int]int)
	}
	m.dirty[key] = value
	m.mu.Unlock()
}

func (m *RWMutexMap) Delete(key int) {
	m.mu.Lock()
	delete(m.dirty, key)
	m.mu.Unlock()
}

// DeepCopyMap is an implementation of mapInterface using a Mutex and
// atomic.Value.  It makes deep copies of the map on every write to avoid
// acquiring the Mutex in Load.
type DeepCopyMap struct {
	mu    sync.Mutex
	clean atomic.Value
}

func (m *DeepCopyMap) Find(key int) (value int, ok bool) {
	clean, _ := m.clean.Load().(map[int]int)
	value, ok = clean[key]
	return value, ok
}

func (m *DeepCopyMap) Insert(key, value int) {
	m.mu.Lock()
	dirty := m.dirty()
	dirty[key] = value
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *DeepCopyMap) Delete(key int) {
	m.mu.Lock()
	dirty := m.dirty()
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *DeepCopyMap) dirty() map[int]int {
	clean, _ := m.clean.Load().(map[int]int)
	dirty := make(map[int]int, len(clean)+1)
	for k, v := range clean {
		dirty[k] = v
	}
	return dirty
}

// sync.Map

type SyncMap[K comparable, V any] struct {
	dirty sync.Map
}

func (m *SyncMap[K, V]) Find(key K) (value V, ok bool) {
	v, ok := m.dirty.Load(key)
	if !ok {
		var v V
		return v, false
	}

	return v.(V), true
}

func (m *SyncMap[K, V]) Insert(key, value int) {
	m.dirty.Store(key, value)
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.dirty.Delete(key)
}
