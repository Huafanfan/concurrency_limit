package limiter

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Huafanfan/concurrency_limit/limits_core"
)

type PartitionedLimiter struct {
	aLimiter *AbstractLimiter

	Partitions        map[string]*Partition
	UnknownPartition  *Partition
	DelayedThreads    atomic.Int64
	MaxDelayedThreads int64
	mu                sync.Mutex
}

func NewPartitionedLimiter(limit limits_core.Limit, partitions []*Partition) *PartitionedLimiter {
	a := &PartitionedLimiter{}
	a.aLimiter = NewAbstractLimiter(limit)
	a.Partitions = make(map[string]*Partition)
	for _, partition := range partitions {
		a.Partitions[partition.Name] = partition
	}
	a.UnknownPartition = NewPartition("unknown")
	a.OnNewLimit(a.aLimiter.GetLimit())
	return a
}
func (pLimiter *PartitionedLimiter) ResolvePartition(ctx context.Context) *Partition {
	value, ok := ctx.Value("name").(string)
	if ok {
		p := pLimiter.Partitions[value]
		if p == nil {
			return pLimiter.UnknownPartition
		}
		return p
	} else {
		return pLimiter.UnknownPartition
	}
}

func (pLimiter *PartitionedLimiter) Acquire(ctx context.Context) limits_core.Listener {
	partition := pLimiter.ResolvePartition(ctx)
	pLimiter.mu.Lock()

	if pLimiter.aLimiter.GetInflight() >= pLimiter.aLimiter.GetLimit() && partition.IsLimitExceeded() {
		pLimiter.mu.Unlock()
		if partition.BackoffMillis > 0 && pLimiter.DelayedThreads.Load() < pLimiter.MaxDelayedThreads {
			pLimiter.DelayedThreads.Add(1)
			time.Sleep(time.Duration(pLimiter.DelayedThreads.Load()))
			pLimiter.DelayedThreads.Add(-1)
		}
		return pLimiter.aLimiter.CreateRejectedListener()
	}
	partition.Acquire()
	pLimiter.mu.Unlock()
	listener := pLimiter.aLimiter.CreateListener()
	return NewPartitionedListener(listener, pLimiter, partition)
}

func (pLimiter *PartitionedLimiter) ReleasePartition(partition *Partition) {
	pLimiter.mu.Lock()
	defer pLimiter.mu.Unlock()

	partition.Release()
}

func (pLimiter *PartitionedLimiter) OnNewLimit(newLimit int) {
	pLimiter.aLimiter.OnNewLimit(newLimit)
	for _, partition := range pLimiter.Partitions {
		partition.UpdateLimit(newLimit)
	}
}

type PartitionedListener struct {
	pLimiter  *PartitionedLimiter
	delegate  limits_core.Listener
	partition *Partition
}

func NewPartitionedListener(delegate limits_core.Listener, pl *PartitionedLimiter, partition *Partition) limits_core.Listener {
	a := &PartitionedListener{
		delegate:  delegate,
		pLimiter:  pl,
		partition: partition,
	}
	return a
}

func (pListener *PartitionedListener) OnSuccess(endTime int64) {
	pListener.delegate.OnSuccess(endTime)
	pListener.pLimiter.ReleasePartition(pListener.partition)
	pListener.UpdateLimit()
}

func (pListener *PartitionedListener) OnIgnore(endTime int64) {
	pListener.delegate.OnIgnore(endTime)
	pListener.pLimiter.ReleasePartition(pListener.partition)
	pListener.UpdateLimit()
}

func (pListener *PartitionedListener) OnDropped(endTime int64) {
	pListener.delegate.OnDropped(endTime)
	pListener.pLimiter.ReleasePartition(pListener.partition)
	pListener.UpdateLimit()
}

func (pListener *PartitionedListener) UpdateLimit() {
	if pListener.pLimiter.aLimiter.GetLimit() != pListener.pLimiter.aLimiter.limitAlgorithm.GetLimit() {
		pListener.pLimiter.OnNewLimit(pListener.pLimiter.aLimiter.limitAlgorithm.GetLimit())
	}
}

type Partition struct {
	Name          string
	Percent       float64
	Limit         int
	Busy          int
	BackoffMillis int64
}

func NewPartition(name string) *Partition {
	return &Partition{Name: name}
}

func (p *Partition) UpdateLimit(totalLimit int) {
	p.Limit = int(math.Max(1, math.Ceil(float64(totalLimit)*p.Percent)))
}

func (p *Partition) IsLimitExceeded() bool {
	return p.Busy >= p.Limit
}

func (p *Partition) Acquire() {
	p.Busy++
}

func (p *Partition) Release() {
	p.Busy--
}
