
# Event-driven and fast I/O multiplexing

Go standard `select` cannot listen on N channels.

## Implementation

"Merging" approache: https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement

## API

### Named channel

- pool := NewNamedPool()
- pool.AddChannel(id, chan)
- pool.RemoveChannel(id)
- pool.Select(func (channelID string, msg interface{}) { ... })

### Unamed channel

- pool := NewPool()
- pool.AddChannel(chan)
- pool.RemoveChannel(chan)
- pool.Select(func (msg interface{}) { ... })

## TODO before library release

- thread safe
- improve API
- unit tests
- doc