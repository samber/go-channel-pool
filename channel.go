package channelPool

type Channel struct {
	ChannelID string
	Chan      chan interface{}
}
