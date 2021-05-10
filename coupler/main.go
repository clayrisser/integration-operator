/**
 * inspired by https://mrkaran.dev/posts/job-queue-golang
 */

package coupler

import (
	"context"
	"math"
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
	OnBroken     func(data interface{})
	OnChanged    func(data interface{})
	OnCreated    func(data interface{})
	OnDeparted   func(data interface{})
	OnJoined     func(data interface{})
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

var GlobalCoupler Coupler = *NewCoupler(Options{
	MaxQueueSize: 99,
	MaxWorkers:   1,
})

func (c *Coupler) Configure(options Options) {
	c.Options = options
}

func (c *Coupler) Start() {
	wg := sync.WaitGroup{}
	maxWorkers := int(math.Max(float64(c.Options.MaxWorkers), 1))
	maxQueueSize := int(math.Max(float64(c.Options.MaxQueueSize), 1))
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		brokenCh := make(chan Event, maxQueueSize)
		changedCh := make(chan Event, maxQueueSize)
		createdCh := make(chan Event, maxQueueSize)
		departedCh := make(chan Event, maxQueueSize)
		joinedCh := make(chan Event, maxQueueSize)
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
					if c.Options.OnCreated != nil {
						c.Options.OnCreated(event.Data)
					}
				case event = <-joinedCh:
					if c.Options.OnJoined != nil {
						c.Options.OnJoined(event.Data)
					}
				case event = <-changedCh:
					if c.Options.OnChanged != nil {
						c.Options.OnChanged(event.Data)
					}
				case event = <-departedCh:
					if c.Options.OnDeparted != nil {
						c.Options.OnDeparted(event.Data)
					}
				case event = <-brokenCh:
					if c.Options.OnBroken != nil {
						c.Options.OnBroken(event.Data)
					}
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
