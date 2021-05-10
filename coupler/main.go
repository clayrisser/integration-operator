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
	"time"
)

type Coupler struct {
	Options  Options
	Finished bool
	bus      *Bus
	closeCh  chan os.Signal
	ctx      context.Context
	cancel   context.CancelFunc
}

type Options struct {
	MaxQueueSize int
	MaxWorkers   int
}

func NewCoupler(options Options) *Coupler {
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
		brokenCh := make(chan Event, c.Options.MaxQueueSize)
		changedCh := make(chan Event, c.Options.MaxQueueSize)
		createdCh := make(chan Event, c.Options.MaxQueueSize)
		departedCh := make(chan Event, c.Options.MaxQueueSize)
		joinedCh := make(chan Event, c.Options.MaxQueueSize)
		c.bus.Sub(BrokenTopic, brokenCh)
		c.bus.Sub(ChangedTopic, changedCh)
		c.bus.Sub(CreatedTopic, createdCh)
		c.bus.Sub(DepartedTopic, departedCh)
		c.bus.Sub(JoinedTopic, joinedCh)
		go func() {
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
		time.Sleep(1 * time.Second)
	}
}

func (c *Coupler) Stop() {
	c.cancel()
}

func (c *Coupler) Wait() {
	<-c.closeCh
	c.Stop()
}

func (c *Coupler) Created(data interface{}) {
	c.bus.Pub(CreatedTopic, data)
}

func (c *Coupler) Joined(data interface{}) {
	c.bus.Pub(JoinedTopic, data)
}

func (c *Coupler) Changed(data interface{}) {
	c.bus.Pub(ChangedTopic, data)
}

func (c *Coupler) Departed(data interface{}) {
	c.bus.Pub(DepartedTopic, data)
}

func (c *Coupler) Broken(data interface{}) {
	c.bus.Pub(BrokenTopic, data)
}
