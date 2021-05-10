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
	Finished bool
	Options  Options
	bus      *Bus
	cancel   context.CancelFunc
	closeCh  chan os.Signal
	ctx      context.Context
	events   *Events
}

type Options struct {
	MaxQueueSize int
	MaxWorkers   int
}

type Events struct {
	OnBroken        func(data interface{})
	OnDeparted      func(data interface{})
	OnJoined        func(data interface{})
	OnPlugChanged   func(data interface{})
	OnPlugCreated   func(data interface{})
	OnSocketChanged func(data interface{})
	OnSocketCreated func(data interface{})
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

func (c *Coupler) Configure(options Options) {
	c.Options = options
}

func (c *Coupler) RegisterEvents(events *Events) {
	c.events = events
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
					if event.Kind == PlugKind {
						if c.events.OnPlugCreated != nil {
							c.events.OnPlugCreated(event.Data)
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketCreated != nil {
							c.events.OnSocketCreated(event.Data)
						}
					}
				case event = <-joinedCh:
					if c.events.OnJoined != nil {
						c.events.OnJoined(event.Data)
					}
				case event = <-changedCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugChanged != nil {
							c.events.OnPlugChanged(event.Data)
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketChanged != nil {
							c.events.OnSocketChanged(event.Data)
						}
					}
				case event = <-departedCh:
					if c.events.OnDeparted != nil {
						c.events.OnDeparted(event.Data)
					}
				case event = <-brokenCh:
					if c.events.OnBroken != nil {
						c.events.OnBroken(event.Data)
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

func (c *Coupler) CreatedPlug(data interface{}) {
	c.bus.Pub(CreatedTopic, PlugKind, data)
}

func (c *Coupler) Joined(data interface{}) {
	c.bus.Pub(JoinedTopic, 0, data)
}

func (c *Coupler) ChangedPlug(data interface{}) {
	c.bus.Pub(ChangedTopic, PlugKind, data)
}

func (c *Coupler) CreatedSocket(data interface{}) {
	c.bus.Pub(CreatedTopic, SocketKind, data)
}

func (c *Coupler) ChangedSocket(data interface{}) {
	c.bus.Pub(ChangedTopic, SocketKind, data)
}

func (c *Coupler) Departed(data interface{}) {
	c.bus.Pub(DepartedTopic, 0, data)
}

func (c *Coupler) Broken(data interface{}) {
	c.bus.Pub(BrokenTopic, 0, data)
}
