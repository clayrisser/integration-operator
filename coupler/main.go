/**
 * inspired by https://mrkaran.dev/posts/job-queue-golang
 */

package coupler

import (
	"context"
	"encoding/json"
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
	OnPlugBroken      func(data interface{}) error
	OnPlugCoupled     func(data interface{}) error
	OnPlugCreated     func(data interface{}) error
	OnPlugDecoupled   func(data interface{}) error
	OnPlugDeleted     func(data interface{}) error
	OnPlugUpdated     func(data interface{}) error
	OnSocketBroken    func(data interface{}) error
	OnSocketCoupled   func(data interface{}) error
	OnSocketCreated   func(data interface{}) error
	OnSocketDecoupled func(data interface{}) error
	OnSocketDeleted   func(data interface{}) error
	OnSocketUpdated   func(data interface{}) error
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
		brokenEventCh := make(chan Event, maxQueueSize)
		updatedEventCh := make(chan Event, maxQueueSize)
		createdEventCh := make(chan Event, maxQueueSize)
		decoupledEventCh := make(chan Event, maxQueueSize)
		deletedEventCh := make(chan Event, maxQueueSize)
		coupledEventCh := make(chan Event, maxQueueSize)
		c.bus.Sub(BrokenTopic, brokenEventCh)
		c.bus.Sub(CoupledTopic, coupledEventCh)
		c.bus.Sub(CreatedTopic, createdEventCh)
		c.bus.Sub(DecoupledTopic, decoupledEventCh)
		c.bus.Sub(DeletedTopic, deletedEventCh)
		c.bus.Sub(UpdatedTopic, updatedEventCh)
		go func() {
			for {
				event := Event{}
				select {
				case <-c.ctx.Done():
					c.bus.Teardown()
					close(brokenEventCh)
					close(updatedEventCh)
					close(createdEventCh)
					close(decoupledEventCh)
					close(coupledEventCh)
					c.Finished = true
					wg.Done()
				case event = <-createdEventCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugCreated != nil {
							if err := c.events.OnPlugCreated(event.Data); err != nil {
								*event.ErrCh <- err
								continue
							}
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketCreated != nil {
							if err := c.events.OnSocketCreated(event.Data); err != nil {
								*event.ErrCh <- err
								continue
							}
						}
					}
					*event.ErrCh <- nil
					continue
				case event = <-coupledEventCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugCoupled != nil {
							if err := c.events.OnPlugCoupled(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketCoupled != nil {
							if err := c.events.OnSocketCoupled(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					}
					*event.ErrCh <- nil
					continue
				case event = <-updatedEventCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugUpdated != nil {
							if err := c.events.OnPlugUpdated(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketUpdated != nil {
							if err := c.events.OnSocketUpdated(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					}
					*event.ErrCh <- nil
					continue
				case event = <-decoupledEventCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugDecoupled != nil {
							if err := c.events.OnPlugDecoupled(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketDecoupled != nil {
							if err := c.events.OnSocketDecoupled(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					}
					*event.ErrCh <- nil
					continue
				case event = <-deletedEventCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugDeleted != nil {
							if err := c.events.OnPlugDeleted(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketDeleted != nil {
							if err := c.events.OnSocketDeleted(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					}
					*event.ErrCh <- nil
					continue
				case event = <-brokenEventCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugBroken != nil {
							if err := c.events.OnPlugBroken(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketBroken != nil {
							if err := c.events.OnSocketBroken(event.Data); err != nil {
								*event.ErrCh <- nil
								continue
							}
						}
					}
					*event.ErrCh <- nil
					continue
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

func (c *Coupler) CreatedPlug(plug *integrationv1alpha2.Plug) error {
	errCh := make(chan error)
	b, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	c.bus.Pub(CreatedTopic, PlugKind, struct {
		ctx  *context.Context
		plug []byte
	}{
		ctx:  &c.ctx,
		plug: b,
	}, errCh)
	return <-errCh
}

func (c *Coupler) CoupledPlug(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig Config,
	socketConfig Config,
) error {
	errCh := make(chan error)
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	c.bus.Pub(CoupledTopic, PlugKind, struct {
		ctx          *context.Context
		plug         []byte
		socket       []byte
		plugConfig   map[string]string
		socketConfig map[string]string
	}{
		ctx:          &c.ctx,
		plug:         bPlug,
		socket:       bSocket,
		plugConfig:   plugConfig,
		socketConfig: socketConfig,
	}, errCh)
	return <-errCh
}

func (c *Coupler) UpdatedPlug(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig Config,
	socketConfig Config,
) error {
	errCh := make(chan error)
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	c.bus.Pub(UpdatedTopic, PlugKind, struct {
		ctx          *context.Context
		plug         []byte
		socket       []byte
		plugConfig   map[string]string
		socketConfig map[string]string
	}{
		ctx:          &c.ctx,
		plug:         bPlug,
		socket:       bSocket,
		plugConfig:   plugConfig,
		socketConfig: socketConfig,
	}, errCh)
	return <-errCh
}

func (c *Coupler) DecoupledPlug(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig Config,
	socketConfig Config,
) error {
	errCh := make(chan error)
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	c.bus.Pub(DecoupledTopic, PlugKind, struct {
		ctx          *context.Context
		plug         []byte
		socket       []byte
		plugConfig   map[string]string
		socketConfig map[string]string
	}{
		ctx:          &c.ctx,
		plug:         bPlug,
		socket:       bSocket,
		plugConfig:   plugConfig,
		socketConfig: socketConfig,
	}, errCh)
	return <-errCh
}

func (c *Coupler) DeletedPlug(
	plug *integrationv1alpha2.Plug,
) error {
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	c.bus.Pub(DeletedTopic, PlugKind, struct {
		ctx  *context.Context
		plug []byte
	}{ctx: &c.ctx, plug: bPlug}, errCh)
	return <-errCh
}

func (c *Coupler) BrokenPlug(
	plug *integrationv1alpha2.Plug,
) error {
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	c.bus.Pub(BrokenTopic, PlugKind, struct {
		ctx  *context.Context
		plug []byte
	}{ctx: &c.ctx, plug: bPlug}, errCh)
	return <-errCh
}

func (c *Coupler) CreatedSocket(
	socket *integrationv1alpha2.Socket,
) error {
	errCh := make(chan error)
	b, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	c.bus.Pub(CreatedTopic, SocketKind, struct {
		ctx    *context.Context
		socket []byte
	}{ctx: &c.ctx, socket: b}, errCh)
	return <-errCh
}

func (c *Coupler) CoupledSocket(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig Config,
	socketConfig Config,
) error {
	errCh := make(chan error)
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	c.bus.Pub(CoupledTopic, SocketKind, struct {
		ctx          *context.Context
		plug         []byte
		socket       []byte
		plugConfig   map[string]string
		socketConfig map[string]string
	}{
		ctx:          &c.ctx,
		plug:         bPlug,
		socket:       bSocket,
		plugConfig:   plugConfig,
		socketConfig: socketConfig,
	}, errCh)
	return <-errCh
}

func (c *Coupler) UpdatedSocket(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig Config,
	socketConfig Config,
) error {
	errCh := make(chan error)
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	c.bus.Pub(UpdatedTopic, SocketKind, struct {
		ctx          *context.Context
		plug         []byte
		socket       []byte
		plugConfig   map[string]string
		socketConfig map[string]string
	}{
		ctx:          &c.ctx,
		plug:         bPlug,
		socket:       bSocket,
		plugConfig:   plugConfig,
		socketConfig: socketConfig,
	}, errCh)
	return <-errCh
}

func (c *Coupler) DecoupledSocket(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig Config,
	socketConfig Config,
) error {
	errCh := make(chan error)
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	c.bus.Pub(DecoupledTopic, SocketKind, struct {
		ctx          *context.Context
		plug         []byte
		socket       []byte
		plugConfig   map[string]string
		socketConfig map[string]string
	}{
		ctx:          &c.ctx,
		plug:         bPlug,
		socket:       bSocket,
		plugConfig:   plugConfig,
		socketConfig: socketConfig,
	}, errCh)
	return <-errCh
}

func (c *Coupler) DeletedSocket(
	socket *integrationv1alpha2.Socket,
) error {
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	c.bus.Pub(DeletedTopic, SocketKind, struct {
		ctx    *context.Context
		socket []byte
	}{ctx: &c.ctx, socket: bSocket}, errCh)
	return <-errCh
}

func (c *Coupler) BrokenSocket(
	socket *integrationv1alpha2.Socket,
) error {
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	c.bus.Pub(BrokenTopic, SocketKind, struct {
		ctx    *context.Context
		socket []byte
	}{ctx: &c.ctx, socket: bSocket}, errCh)
	return <-errCh
}
