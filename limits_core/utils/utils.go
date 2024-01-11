package utils

import "sync/atomic"

type AtomicLongArray struct {
	data []int64
}

func NewAtomicLongArray(length int) *AtomicLongArray {
	return &AtomicLongArray{
		data: make([]int64, length),
	}
}

func (arr *AtomicLongArray) Get(index int) int64 {
	return atomic.LoadInt64(&arr.data[index])
}

func (arr *AtomicLongArray) Set(index int, value int64) {
	atomic.StoreInt64(&arr.data[index], value)
}

func (arr *AtomicLongArray) Length() int {
	return len(arr.data)
}
