package main

import "sync"

type KV struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
		// mu:   sync.RWMutex{},
	}
}

func (kv *KV) Set(key, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[string(key)] = val
	return nil
}

func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[string(key)]
	return val, ok
}
