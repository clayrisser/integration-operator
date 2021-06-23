/**
 * inspired by https://mrkaran.dev/posts/job-queue-golang
 */

package coupler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	minifyJson "github.com/tdewolff/minify/json"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"

	"github.com/go-resty/resty/v2"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/util"
	"github.com/tdewolff/minify"
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
	c.bus.Pub(CreatedTopic, PlugKind, struct{ plug []byte }{plug: b}, errCh)
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
		plug         []byte
		socket       []byte
		plugConfig   []byte
		socketConfig []byte
	}{plug: bPlug, socket: bSocket, plugConfig: plugConfig, socketConfig: socketConfig}, errCh)
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
		plug         []byte
		socket       []byte
		plugConfig   []byte
		socketConfig []byte
	}{plug: bPlug, socket: bSocket, plugConfig: plugConfig, socketConfig: socketConfig}, errCh)
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
		plug         []byte
		socket       []byte
		plugConfig   []byte
		socketConfig []byte
	}{plug: bPlug, socket: bSocket, plugConfig: plugConfig, socketConfig: socketConfig}, errCh)
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
		plug []byte
	}{plug: bPlug}, errCh)
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
		plug []byte
	}{plug: bPlug}, errCh)
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
		plug         []byte
		socket       []byte
		plugConfig   []byte
		socketConfig []byte
	}{plug: bPlug, socket: bSocket, plugConfig: plugConfig, socketConfig: socketConfig}, errCh)
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
		plug         []byte
		socket       []byte
		plugConfig   []byte
		socketConfig []byte
	}{plug: bPlug, socket: bSocket, plugConfig: plugConfig, socketConfig: socketConfig}, errCh)
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
		plug         []byte
		socket       []byte
		plugConfig   []byte
		socketConfig []byte
	}{plug: bPlug, socket: bSocket, plugConfig: plugConfig, socketConfig: socketConfig}, errCh)
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
		socket []byte
	}{socket: bSocket}, errCh)
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
		socket []byte
	}{socket: bSocket}, errCh)
	return <-errCh
}

func (c *Coupler) GetConfig(
	endpoint string,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
) (Config, error) {
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return nil, err
	}
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return nil, err
	}
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	m := minify.New()
	m.AddFunc("application/json", minifyJson.Minify)
	go func() {
		body := `{"version":"1"`
		if plug != nil {
			jsonPlug := gjson.Parse(string(bPlug))
			body += fmt.Sprintf(`,"plug":%s`, jsonPlug)
			meta, _ := yaml.YAMLToJSON([]byte(jsonPlug.Get("spec").Get("meta").String()))
			if meta == nil {
				meta = []byte("{}")
			}
			body += fmt.Sprintf(`,"plugMeta":%s`, meta)
		}
		if socket != nil {
			jsonSocket := gjson.Parse(string(bSocket))
			body += fmt.Sprintf(`,"socket":%s`, jsonSocket)
			meta, _ := yaml.YAMLToJSON([]byte(jsonSocket.Get("spec").Get("meta").String()))
			if meta == nil {
				meta = []byte("{}")
			}
			body += fmt.Sprintf(`,"socketMeta":%s`, meta)
		}
		body += "}"
		body, err := m.String("application/json", body)
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(util.GetEndpoint(endpoint) + "/config")
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
	}
}
