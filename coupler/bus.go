/**
 * inspired by https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997
 */

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
	Kind  Kind
	Data  interface{}
	Topic Topic
}

type Bus struct {
	subscribers map[Topic][](chan<- Event)
	rm          sync.RWMutex
	closed      bool
}

func NewBus() *Bus {
	return &Bus{
		closed:      false,
		rm:          sync.RWMutex{},
		subscribers: map[Topic][](chan<- Event){},
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
	b.subscribers = map[Topic][](chan<- Event){}
}

func (b *Bus) Pub(topic Topic, kind Kind, data interface{}) {
	b.rm.RLock()
	if channels, found := b.subscribers[topic]; found {
		go func(event Event, channels [](chan<- Event)) {
			for _, ch := range channels {
				if !b.closed {
					ch <- event
				}
			}
		}(Event{Topic: topic, Data: data, Kind: kind}, append([](chan<- Event){}, channels...))
	}
	b.rm.RUnlock()
}

func (b *Bus) Sub(topic Topic, ch chan<- Event) {
	b.rm.Lock()
	if topicSubscribers, found := b.subscribers[topic]; found {
		b.subscribers[topic] = append(topicSubscribers, ch)
	} else {
		b.subscribers[topic] = append([](chan<- Event){}, ch)
	}
	b.rm.Unlock()
}
