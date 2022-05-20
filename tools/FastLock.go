package tools

import (
	"runtime"
	"sync/atomic"
)

const (
	UNLOCKED = 0
	LOCKED   = 1
)

const (
	MAX_LOOP_COUNT = 100000
)

type FastLock struct {
	flag      int32
	loopCount int32
	name      string
}

func NewFastLock(name string) *FastLock {
	this := new(FastLock)
	this.flag = UNLOCKED
	this.loopCount = MAX_LOOP_COUNT
	this.name = name
	return this
}

func NewFastLockWithLoopCount(loopCount int32) *FastLock {
	this := new(FastLock)
	this.flag = UNLOCKED
	if loopCount <= 0 {
		loopCount = MAX_LOOP_COUNT
	}
	this.loopCount = loopCount
	return this
}

func (this *FastLock) Lock() bool {

	for i := int32(0); i < this.loopCount; i++ {
		if atomic.CompareAndSwapInt32(&this.flag, UNLOCKED, LOCKED) {
			return true
		}
		runtime.Gosched()
	}

	ShowWarnning("FastLock.Lock failed:", this.name)
	return false
}

func (this *FastLock) Unlock() bool {
	for i := int32(0); i < this.loopCount; i++ {
		if atomic.CompareAndSwapInt32(&this.flag, LOCKED, UNLOCKED) || atomic.LoadInt32(&this.flag) == UNLOCKED {
			return true
		}
		runtime.Gosched()
	}
	ShowWarnning("FastLock.Unlock failed:", this.name)
	return false
}

func (this *FastLock) UnlockForce() {
	for {
		if atomic.CompareAndSwapInt32(&this.flag, LOCKED, UNLOCKED) || atomic.LoadInt32(&this.flag) == UNLOCKED {
			return
		}
		runtime.Gosched()
	}
}
