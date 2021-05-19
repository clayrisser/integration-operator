package util

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Scope int

const (
	UpdateScope Scope = iota + 1
	UpdateStatusScope
)

const (
	UpdatePlugTopic Topic = iota + 1
	UpdateSocketTopic
)

type Update struct {
	Finished     bool
	MaxQueueSize int
	bus          *Bus
	cancel       context.CancelFunc
	closeCh      chan os.Signal
	ctx          context.Context
}

func NewUpdate(maxQueueSize int) *Update {
	if maxQueueSize == 0 {
		maxQueueSize = 99
	}
	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	return &Update{
		Finished:     false,
		MaxQueueSize: maxQueueSize,
		bus:          NewBus(),
		cancel:       cancel,
		closeCh:      closeCh,
		ctx:          ctx,
	}
}

func (u *Update) SchedulePlugUpdate(
	c *client.Client,
	ctx *context.Context,
	log *logr.Logger,
	plug *integrationv1alpha2.Plug,
) {
	u.bus.Pub(UpdatePlugTopic,
		struct{ Scope Scope }{Scope: UpdateScope},
		struct {
			Client *client.Client
			Ctx    *context.Context
			Log    *logr.Logger
			Plug   *integrationv1alpha2.Plug
		}{
			Client: c,
			Ctx:    ctx,
			Log:    log,
			Plug:   plug,
		},
	)
}

func (u *Update) SchedulePlugUpdateStatus(
	c *client.Client,
	ctx *context.Context,
	log *logr.Logger,
	plug *integrationv1alpha2.Plug,
) {
	u.bus.Pub(UpdatePlugTopic,
		struct{ Scope Scope }{Scope: UpdateStatusScope},
		struct {
			Client *client.Client
			Ctx    *context.Context
			Log    *logr.Logger
			Plug   *integrationv1alpha2.Plug
		}{
			Client: c,
			Ctx:    ctx,
			Log:    log,
			Plug:   plug,
		},
	)
}

func (u *Update) ScheduleSocketUpdate(
	c *client.Client,
	ctx *context.Context,
	log *logr.Logger,
	plug *integrationv1alpha2.Socket,
) {
	u.bus.Pub(UpdateSocketTopic,
		struct{ Scope Scope }{Scope: UpdateScope},
		struct {
			Client *client.Client
			Ctx    *context.Context
			Log    *logr.Logger
			Socket *integrationv1alpha2.Socket
		}{
			Client: c,
			Ctx:    ctx,
			Log:    log,
			Socket: plug,
		},
	)
}

func (u *Update) ScheduleSocketUpdateStatus(
	c *client.Client,
	ctx *context.Context,
	log *logr.Logger,
	plug *integrationv1alpha2.Socket,
) {
	u.bus.Pub(UpdateSocketTopic,
		struct{ Scope Scope }{Scope: UpdateStatusScope},
		struct {
			Client *client.Client
			Ctx    *context.Context
			Log    *logr.Logger
			Socket *integrationv1alpha2.Socket
		}{
			Client: c,
			Ctx:    ctx,
			Log:    log,
			Socket: plug,
		},
	)
}

func (u *Update) Start() {
	wg := sync.WaitGroup{}
	plugCh := make(chan Event, u.MaxQueueSize)
	socketCh := make(chan Event, u.MaxQueueSize)
	u.bus.Sub(UpdatePlugTopic, plugCh)
	u.bus.Sub(UpdateSocketTopic, socketCh)
	wg.Add(1)
	go func() {
		for {
			event := Event{}
			select {
			case <-u.ctx.Done():
				u.bus.Teardown()
				close(plugCh)
				close(socketCh)
				u.Finished = true
				wg.Done()
			case event = <-plugCh:
				d := event.Data.(struct {
					Client *client.Client
					Ctx    *context.Context
					Log    *logr.Logger
					Plug   *integrationv1alpha2.Plug
				})
				m := event.Meta.(struct{ Scope Scope })
				client := *d.Client
				ctx := *d.Ctx
				log := *d.Log
				plug := d.Plug.DeepCopy()
				scope := m.Scope
				if scope == UpdateStatusScope {
					if err := client.Status().Update(ctx, plug); err != nil {
						log.Error(err, "failed to update plug status")
					}
				} else if scope == UpdateScope {
					if err := client.Update(ctx, plug); err != nil {
						log.Error(err, "failed to update plug")
					}
				}
			case event = <-socketCh:
				d := event.Data.(struct {
					client *client.Client
					ctx    *context.Context
					log    *logr.Logger
					socket *integrationv1alpha2.Socket
				})
				m := event.Meta.(struct{ Scope Scope })
				client := *d.client
				ctx := *d.ctx
				log := *d.log
				scope := m.Scope
				socket := d.socket.DeepCopy()
				if scope == UpdateStatusScope {
					if err := client.Status().Update(ctx, socket); err != nil {
						log.Error(err, "failed to update socket status")
					}
				} else if scope == UpdateScope {
					if err := client.Update(ctx, socket); err != nil {
						log.Error(err, "failed to update socket")
					}
				}
			}
		}
	}()
}

func (u *Update) Stop() {
	u.cancel()
}

func (u *Update) Wait() {
	<-u.closeCh
	u.Stop()
}

func (u *Update) Run() {
	u.Start()
	u.Wait()
}
