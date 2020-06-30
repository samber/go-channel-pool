package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	pooling "channel-pool"

	randomdata "github.com/Pallinder/go-randomdata"
	uuid "github.com/satori/go.uuid"
)

var done chan os.Signal

func main() {
	pool := pooling.NewNamedPool()
	defer pool.Close()

	produceRandomData(pool)
	handleStop(pool)

	// listen on all channels
	for {
		pool.Select(func(channelID string, msg interface{}) {
			str, ok := msg.(string)
			if ok {
				fmt.Printf("Message from %s channel: %s\n", channelID, str)
			}
		}, func(channelID string) {
			fmt.Printf("Closed channel %s\n", channelID)
		})
	}
}

func produceRandomData(pool *pooling.NamedPool) {
	setInterval(func() {
		createOrRemoveChannels(pool)
	}, 500*time.Millisecond)
	setInterval(func() {
		produceMessage(pool)
	}, 100*time.Millisecond)
}

// kill on ctrl-c
func handleStop(pool *pooling.NamedPool) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		fmt.Println("Stopping")
		for _, c := range pool.GetChannels() {
			close(c.Chan)
		}
	}()
}

func setInterval(f func(), t time.Duration) {
	ticker := time.NewTicker(t)
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				f()
			}
		}
	}()
}

// randomly add or remove some channels
func createOrRemoveChannels(pool *pooling.NamedPool) {
	now := time.Now().Nanosecond()
	if now/1000%5 == 0 {
		pool.AddChannel(uuid.NewV4().String(), make(chan interface{}))
		fmt.Printf("%d channel(s)\n", len(pool.GetChannels()))
	} else if now/1000%5 == 1 {
		channels := pool.GetChannels()
		if len(channels) > 0 {
			c := channels[rand.Intn(len(channels))]
			pool.RemoveChannel(c.ChannelID)
			close(c.Chan)
			fmt.Printf("%d channel(s)\n", len(pool.GetChannels()))
		}
	}
}

// send country name randomly
func produceMessage(pool *pooling.NamedPool) {
	channels := pool.GetChannels()
	if len(channels) == 0 {
		return
	}

	defer func() {
		// in case of race condition
		recover()
	}()
	channels[rand.Intn(len(channels))].Chan <- randomdata.Country(randomdata.FullCountry)
}
