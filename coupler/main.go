/**
 * inspired by https://mrkaran.dev/posts/job-queue-golang
 */

package coupler

import (
	"context"
	"fmt"
	"sync"
)

type Coupler struct {
	Options  CouplerOptions
	Finished bool
	bus      *Bus
}

type CouplerOptions struct {
	MaxQueueSize int
	MaxWorkers   int
}

func NewCoupler(options CouplerOptions) *Coupler {
	return &Coupler{
		Finished: false,
		Options:  options,
		bus:      NewBus(),
	}
}

func (c *Coupler) Run(ctx context.Context) {
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
				case <-ctx.Done():
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
