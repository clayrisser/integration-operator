/**
 * inspired by https://mrkaran.dev/posts/job-queue-golang
 */

package coupler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Coupler struct {
	Options  CouplerOptions
	Finished bool
	bus      *Bus
	closeCh  chan os.Signal
	ctx      context.Context
	cancel   context.CancelFunc
}

type CouplerOptions struct {
	MaxQueueSize int
	MaxWorkers   int
}

func NewCoupler(options CouplerOptions) *Coupler {
	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	return &Coupler{
		Finished: false,
		Options:  options,
		bus:      NewBus(),
		cancel:   cancel,
		closeCh:  closeCh,
		ctx:      ctx,
	}
}

func (c *Coupler) Start() {
	wg := sync.WaitGroup{}
	for i := 0; i < c.Options.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			brokenCh := make(chan Event, c.Options.MaxQueueSize)
			changedCh := make(chan Event, c.Options.MaxQueueSize)
			createdCh := make(chan Event, c.Options.MaxQueueSize)
			departedCh := make(chan Event, c.Options.MaxQueueSize)
			joinedCh := make(chan Event, c.Options.MaxQueueSize)
			c.bus.Sub(Broken, brokenCh)
			c.bus.Sub(Changed, changedCh)
			c.bus.Sub(Created, createdCh)
			c.bus.Sub(Departed, departedCh)
			c.bus.Sub(Joined, joinedCh)
			for {
				event := Event{}
				select {
				case <-c.ctx.Done():
					c.bus.Teardown()
					close(brokenCh)
					close(changedCh)
					close(createdCh)
					close(departedCh)
					close(joinedCh)
					c.Finished = true
					wg.Done()
				case event = <-createdCh:
					fmt.Println(event.Data)
					fmt.Println(event.Topic)
				case event = <-joinedCh:
					fmt.Println(event.Data)
					fmt.Println(event.Topic)
				case event = <-changedCh:
					fmt.Println(event.Data)
					fmt.Println(event.Topic)
				case event = <-departedCh:
					fmt.Println(event.Data)
					fmt.Println(event.Topic)
				case event = <-brokenCh:
					fmt.Println(event.Data)
					fmt.Println(event.Topic)
				}
			}
		}()
	}
}

func (c *Coupler) Stop() {
	c.cancel()
}

func (c *Coupler) Wait() {
	<-c.closeCh
	c.Stop()
}
