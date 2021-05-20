/**
 * inspired by https://mrkaran.dev/posts/job-queue-golang
 */

package coupler

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
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
	OnBroken        func(data interface{}) error
	OnDecoupled     func(data interface{}) error
	OnPlugCoupled   func(data interface{}) error
	OnPlugCreated   func(data interface{}) error
	OnPlugUpdated   func(data interface{}) error
	OnSocketCoupled func(data interface{}) error
	OnSocketCreated func(data interface{}) error
	OnSocketUpdated func(data interface{}) error
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
		coupledEventCh := make(chan Event, maxQueueSize)
		c.bus.Sub(BrokenTopic, brokenEventCh)
		c.bus.Sub(UpdatedTopic, updatedEventCh)
		c.bus.Sub(CreatedTopic, createdEventCh)
		c.bus.Sub(DecoupledTopic, decoupledEventCh)
		c.bus.Sub(CoupledTopic, coupledEventCh)
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
					if c.events.OnDecoupled != nil {
						if err := c.events.OnDecoupled(event.Data); err != nil {
							*event.ErrCh <- nil
							continue
						}
					}
					*event.ErrCh <- nil
					continue
				case event = <-brokenEventCh:
					if c.events.OnBroken != nil {
						if err := c.events.OnBroken(event.Data); err != nil {
							*event.ErrCh <- nil
							continue
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
	c.bus.Pub(CreatedTopic, PlugKind, struct{ plug []byte }{plug: b}, errCh)
	return <-errCh
}

func (c *Coupler) CoupledPlug(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	config Config,
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
		plug   []byte
		socket []byte
		config []byte
	}{plug: bPlug, socket: bSocket, config: config}, errCh)
	return <-errCh
}

func (c *Coupler) UpdatedPlug(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	config Config,
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
		plug   []byte
		socket []byte
		config []byte
	}{plug: bPlug, socket: bSocket, config: config}, errCh)
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
	c.bus.Pub(CreatedTopic, SocketKind, struct{ socket []byte }{socket: b}, errCh)
	return <-errCh
}

func (c *Coupler) CoupledSocket(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	config Config,
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
		plug   []byte
		socket []byte
		config []byte
	}{plug: bPlug, socket: bSocket, config: config}, errCh)
	return <-errCh
}

func (c *Coupler) UpdatedSocket(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	config Config,
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
		plug   []byte
		socket []byte
		config []byte
	}{plug: bPlug, socket: bSocket, config: config}, errCh)
	return <-errCh
}

func (c *Coupler) Decoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
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
	c.bus.Pub(DecoupledTopic, 0, struct {
		plug   []byte
		socket []byte
	}{plug: bPlug, socket: bSocket}, errCh)
	return <-errCh
}

func (c *Coupler) Broken() error {
	errCh := make(chan error)
	c.bus.Pub(BrokenTopic, 0, struct{}{}, errCh)
	return <-errCh
}

func (c *Coupler) GetConfig(endpoint string) (Config, error) {
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	go func() {
		r, err := client.R().EnableTrace().SetQueryParams(map[string]string{
			"version": "1",
		}).Get(endpoint)
		if err != nil {
			errCh <- err
		}
		rCh <- r
	}()
	select {
	case r := <-rCh:
		return r.Body(), nil
	case err := <-errCh:
		return nil, err
	case <-time.After(3 * time.Second):
		return nil, errors.New("timeout")
	}
}
