package services

import (
	"sync"
	"sync/atomic"
)

type ErrorsManager interface {
	Store(errorMsg string) int32
	GetCount(errorMsg string) int32
}

type atomicCounter int32

func (ac *atomicCounter) plusOne() int32 {
	return atomic.AddInt32((*int32)(ac), 1)
}
func (ac *atomicCounter) current() int32 {
	return atomic.LoadInt32((*int32)(ac))
}

type errorsManager struct {
	eCounterMapMutex sync.RWMutex
	eCounterMap      map[string]atomicCounter
}

func NewErrorsManager() ErrorsManager {
	return &errorsManager{
		eCounterMapMutex: sync.RWMutex{},
		eCounterMap:      make(map[string]atomicCounter),
	}
}

func (em *errorsManager) Store(errorMsg string) int32 {
	em.eCounterMapMutex.Lock()
	currVal := em.eCounterMap[errorMsg]
	currVal.plusOne()
	em.eCounterMap[errorMsg] = currVal
	em.eCounterMapMutex.Unlock()
	return currVal.current()
}

func (em *errorsManager) GetCount(errorMsg string) int32 {
	// here the mismatch happens also using RLock/RUnlock
	em.eCounterMapMutex.Lock()
	count := em.eCounterMap[errorMsg]
	em.eCounterMapMutex.Unlock()
	return count.current()
}
