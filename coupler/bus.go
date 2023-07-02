/**
 * File: /bus.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 02-07-2023 11:49:19
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

// inspired by https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997

package coupler

import "sync"

type Topic int

const (
	CreatedTopic Topic = iota + 1
	CoupledTopic
	UpdatedTopic
	DecoupledTopic
	DeletedTopic
	BrokenTopic
)

type Kind int

const (
	PlugKind Kind = iota + 1
	SocketKind
)

type Event struct {
	Data  interface{}
	ErrCh *chan error
	Kind  Kind
	Topic Topic
}

type Bus struct {
	subscribers map[Topic][]*chan<- Event
	rm          sync.RWMutex
	closed      bool
}

func NewBus() *Bus {
	return &Bus{
		closed:      false,
		rm:          sync.RWMutex{},
		subscribers: map[Topic][]*chan<- Event{},
	}
}

func (b *Bus) Close() {
	b.closed = true
}

func (b *Bus) Open() {
	b.closed = false
}

func (b *Bus) Teardown() {
	b.Close()
	b.subscribers = map[Topic][]*chan<- Event{}
}

func (b *Bus) Pub(topic Topic, kind Kind, data interface{}, errCh chan error) {
	b.rm.RLock()
	if eventChannels, found := b.subscribers[topic]; found {
		go func(event Event, eventChannels []*chan<- Event) {
			for _, eventCh := range eventChannels {
				if !b.closed {
					*eventCh <- event
				}
			}
		}(Event{Topic: topic, Data: data, Kind: kind, ErrCh: &errCh}, append([]*chan<- Event{}, eventChannels...))
	}
	b.rm.RUnlock()
}

func (b *Bus) Sub(topic Topic, eventCh chan<- Event) {
	b.rm.Lock()
	if topicSubscribers, found := b.subscribers[topic]; found {
		b.subscribers[topic] = append(topicSubscribers, &eventCh)
	} else {
		b.subscribers[topic] = append([]*chan<- Event{}, &eventCh)
	}
	b.rm.Unlock()
}
