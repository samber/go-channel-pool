package channelPool

import "sync"

type message struct {
	channelID string
	msg       interface{}
}

type NamedPool struct {
	agg      chan message
	channels map[string]*Channel
	mux      sync.Mutex
}

func NewNamedPool() *NamedPool {
	return &NamedPool{
		agg:      make(chan message),
		channels: map[string]*Channel{},
	}
}

func (p *NamedPool) Select(callback func(string, interface{}), closed func(string)) {
	select {
	case msg, done := <-p.agg:
		if done {
			callback(msg.channelID, msg.msg)
		} else {
			closed(msg.channelID)
		}
	}
}

func (p *NamedPool) Close() {
	p.channels = map[string]*Channel{}
	close(p.agg)
}

func (p *NamedPool) AddChannel(id string, in chan interface{}) {
	p.mux.Lock()
	p.removeChannel(id)

	c := &Channel{
		ChannelID: id,
		Chan:      in,
	}

	go p.streamChannel(c)
	p.channels[c.ChannelID] = c

	p.mux.Unlock()
}

func (p *NamedPool) streamChannel(c *Channel) {
	for msg := range c.Chan {
		p.agg <- message{channelID: c.ChannelID, msg: msg}
	}
}

func (p *NamedPool) RemoveChannel(id string) {
	p.mux.Lock()
	p.removeChannel(id)
	p.mux.Unlock()
}

// not thread safe
func (p *NamedPool) removeChannel(id string) {
	if _, ok := p.channels[id]; ok {
		delete(p.channels, id)
		agg := make(chan message)

		old := p.agg
		p.agg = agg
		close(old)
	}
}

func (p *NamedPool) GetChannels() []*Channel {
	channels := []*Channel{}
	for _, v := range p.channels {
		channels = append(channels, v)
	}
	return channels
}
