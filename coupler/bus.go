/*
 * File: /coupler/bus.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 26-06-2021 10:54:56
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Silicon Hills LLC (c) Copyright 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
