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

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
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
	OnPlugBroken     func(data interface{})
	OnPlugChanged    func(data interface{})
	OnPlugCreated    func(data interface{})
	OnPlugDeparted   func(data interface{})
	OnPlugJoined     func(data interface{})
	OnSocketBroken   func(data interface{})
	OnSocketChanged  func(data interface{})
	OnSocketCreated  func(data interface{})
	OnSocketDeparted func(data interface{})
	OnSocketJoined   func(data interface{})
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
					if event.Kind == PlugKind {
						if c.events.OnPlugJoined != nil {
							c.events.OnPlugJoined(event.Data)
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketJoined != nil {
							c.events.OnSocketJoined(event.Data)
						}
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
					if event.Kind == PlugKind {
						if c.events.OnPlugDeparted != nil {
							c.events.OnPlugDeparted(event.Data)
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketDeparted != nil {
							c.events.OnSocketDeparted(event.Data)
						}
					}
				case event = <-brokenCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugBroken != nil {
							c.events.OnPlugBroken(event.Data)
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketBroken != nil {
							c.events.OnSocketBroken(event.Data)
						}
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

func (c *Coupler) JoinedPlug(data interface{}) {
	c.bus.Pub(JoinedTopic, PlugKind, data)
}

func (c *Coupler) ChangedPlug(data interface{}) {
	c.bus.Pub(ChangedTopic, PlugKind, data)
}

func (c *Coupler) DepartedPlug(data interface{}) {
	c.bus.Pub(DepartedTopic, PlugKind, data)
}

func (c *Coupler) BrokenPlug(data interface{}) {
	c.bus.Pub(BrokenTopic, PlugKind, data)
}

func (c *Coupler) CreatedSocket(data interface{}) {
	c.bus.Pub(CreatedTopic, SocketKind, data)
}

func (c *Coupler) JoinedSocket(data interface{}) {
	c.bus.Pub(JoinedTopic, SocketKind, data)
}

func (c *Coupler) ChangedSocket(data interface{}) {
	c.bus.Pub(ChangedTopic, SocketKind, data)
}

func (c *Coupler) DepartedSocket(data interface{}) {
	c.bus.Pub(DepartedTopic, SocketKind, data)
}

func (c *Coupler) BrokenSocket(data interface{}) {
	c.bus.Pub(BrokenTopic, SocketKind, data)
}

func CreateGlobalCoupler() Coupler {
	handlers := NewHandlers()
	globalCoupler := *NewCoupler(Options{
		MaxQueueSize: 99,
		MaxWorkers:   1,
	})
	globalCoupler.RegisterEvents(&Events{
		OnPlugCreated: func(data interface{}) {
			d := data.(struct {
				plug *integrationv1alpha2.Plug
			})
			handlers.HandlePlugCreated(d.plug)
		},
		OnPlugBroken: func(data interface{}) {
			d := data.(struct {
				socket  *integrationv1alpha2.Socket
				plug    *integrationv1alpha2.Plug
				payload interface{}
			})
			handlers.HandlePlugBroken(d.socket, d.plug, d.payload)
		},
	})
	return globalCoupler
}

var GlobalCoupler Coupler = CreateGlobalCoupler()
