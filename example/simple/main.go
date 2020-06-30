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
)

var done chan os.Signal

func main() {
	pool := pooling.NewPool()
	defer pool.Close()

	produceRandomData(pool)
	handleStop(pool)

	// listen on all channels
	for {
		pool.Select(func(msg interface{}) {
			str, ok := msg.(string)
			if ok {
				fmt.Printf("Message: %s\n", str)
			}
		}, func() {
			fmt.Println("Closed channel")
		})
	}
}

func produceRandomData(pool *pooling.Pool) {
	setInterval(func() {
		createOrRemoveChannels(pool)
	}, 500*time.Millisecond)
	setInterval(func() {
		produceMessage(pool)
	}, 100*time.Millisecond)
}

// kill on ctrl-c
func handleStop(pool *pooling.Pool) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		fmt.Println("Stopping")
		for _, c := range pool.GetChannels() {
			close(c)
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
func createOrRemoveChannels(pool *pooling.Pool) {
	now := time.Now().Nanosecond()
	if now/1000%5 == 0 {
		pool.AddChannel(make(chan interface{}))
		fmt.Printf("%d channel(s)\n", len(pool.GetChannels()))
	} else if now/1000%5 == 1 {
		channels := pool.GetChannels()
		if len(channels) > 0 {
			c := channels[rand.Intn(len(channels))]
			pool.RemoveChannel(c)
			close(c)
			fmt.Printf("%d channel(s)\n", len(pool.GetChannels()))
		}
	}
}

// send country name randomly
func produceMessage(pool *pooling.Pool) {
	channels := pool.GetChannels()
	if len(channels) == 0 {
		return
	}

	defer func() {
		// in case of race condition
		recover()
	}()
	channels[rand.Intn(len(channels))] <- randomdata.Country(randomdata.FullCountry)
}
