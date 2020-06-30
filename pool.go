package channelPool

import (
	"sync"
)

type Pool struct {
	agg      chan interface{}
	channels []chan interface{}
	mux      sync.Mutex
}

func NewPool() *Pool {
	return &Pool{
		agg:      make(chan interface{}),
		channels: []chan interface{}{},
	}
}

func (p *Pool) Select(callback func(interface{}), closed func()) {
	select {
	case msg, done := <-p.agg:
		if done {
			callback(msg)
		} else {
			closed()
		}
	}
}

func (p *Pool) Close() {
	p.channels = []chan interface{}{}
	close(p.agg)
}

func (p *Pool) AddChannel(channel chan interface{}) {
	p.mux.Lock()
	p.removeChannel(channel)

	go func() {
		for msg := range channel {
			p.agg <- msg
		}
	}()
	p.channels = append(p.channels, channel)

	p.mux.Unlock()
}

func (p *Pool) RemoveChannel(channel chan interface{}) {
	p.mux.Lock()
	p.removeChannel(channel)
	p.mux.Unlock()
}

// not thread safe
func (p *Pool) removeChannel(channel chan interface{}) {
	index := p.findChannel(channel)
	if index != -1 {
		p.channels = append(p.channels[:index], p.channels[index+1:]...)

		agg := make(chan interface{})
		old := p.agg
		p.agg = agg
		close(old)
	}
}

func (p *Pool) GetChannels() []chan interface{} {
	return p.channels
}

func (p *Pool) findChannel(channel chan interface{}) int {
	for i, c := range p.channels {
		if c == channel {
			return i
		}
	}
	return -1
}
