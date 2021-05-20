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
	OnBroken        func(data interface{})
	OnDecoupled     func(data interface{})
	OnPlugCoupled   func(data interface{})
	OnPlugCreated   func(data interface{})
	OnPlugUpdated   func(data interface{})
	OnSocketCoupled func(data interface{})
	OnSocketCreated func(data interface{})
	OnSocketUpdated func(data interface{})
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
		updatedCh := make(chan Event, maxQueueSize)
		createdCh := make(chan Event, maxQueueSize)
		decoupledCh := make(chan Event, maxQueueSize)
		coupledCh := make(chan Event, maxQueueSize)
		c.bus.Sub(BrokenTopic, brokenCh)
		c.bus.Sub(UpdatedTopic, updatedCh)
		c.bus.Sub(CreatedTopic, createdCh)
		c.bus.Sub(DecoupledTopic, decoupledCh)
		c.bus.Sub(CoupledTopic, coupledCh)
		go func() {
			for {
				event := Event{}
				select {
				case <-c.ctx.Done():
					c.bus.Teardown()
					close(brokenCh)
					close(updatedCh)
					close(createdCh)
					close(decoupledCh)
					close(coupledCh)
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
				case event = <-coupledCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugCoupled != nil {
							c.events.OnPlugCoupled(event.Data)
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketCoupled != nil {
							c.events.OnSocketCoupled(event.Data)
						}
					}
				case event = <-updatedCh:
					if event.Kind == PlugKind {
						if c.events.OnPlugUpdated != nil {
							c.events.OnPlugUpdated(event.Data)
						}
					} else if event.Kind == SocketKind {
						if c.events.OnSocketUpdated != nil {
							c.events.OnSocketUpdated(event.Data)
						}
					}
				case event = <-decoupledCh:
					if c.events.OnDecoupled != nil {
						c.events.OnDecoupled(event.Data)
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

func (c *Coupler) CreatedPlug(plug *integrationv1alpha2.Plug) error {
	b, err := json.Marshal(plug)
	if err != nil {
		return err
	}
	c.bus.Pub(CreatedTopic, PlugKind, struct{ plug []byte }{plug: b})
	return nil
}

func (c *Coupler) CoupledPlug(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	config Config,
) error {
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
	}{plug: bPlug, socket: bSocket, config: config})
	return nil
}

func (c *Coupler) UpdatedPlug(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	config Config,
) error {
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
	}{plug: bPlug, socket: bSocket, config: config})
	return nil
}

func (c *Coupler) CreatedSocket(
	socket *integrationv1alpha2.Socket,
) error {
	b, err := json.Marshal(socket)
	if err != nil {
		return err
	}
	c.bus.Pub(CreatedTopic, SocketKind, struct{ socket []byte }{socket: b})
	return nil
}

func (c *Coupler) CoupledSocket(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	config Config,
) error {
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
	}{plug: bPlug, socket: bSocket, config: config})
	return nil
}

func (c *Coupler) UpdatedSocket(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	config Config,
) error {
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
	}{plug: bPlug, socket: bSocket, config: config})
	return nil
}

func (c *Coupler) Decoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
) error {
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
	}{plug: bPlug, socket: bSocket})
	return nil
}

func (c *Coupler) Broken() {
	c.bus.Pub(BrokenTopic, 0, struct{}{})
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
