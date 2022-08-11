package core

import "sync"

// ChannelPool is a worker pool that limits concurrency using a fixed-size buffered channel.
type ChannelPool struct {
	wip chan struct{}
	wg  *sync.WaitGroup
}

func NewChannelPool(workerCount int) ChannelPool {
	wip := make(chan struct{}, workerCount)
	wg := &sync.WaitGroup{}
	return ChannelPool{wip, wg}
}

func (p ChannelPool) Submit(f func()) {
	p.wg.Add(1)
	go func() {
		p.wip <- struct{}{}
		f()
		<-p.wip
		p.wg.Done()
	}()
}

func (p ChannelPool) Wait() {
	p.wg.Wait()
	close(p.wip)
}
