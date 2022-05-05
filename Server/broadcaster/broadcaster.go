package broadcaster

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Broadcaster struct {
	mu      sync.Mutex
	clients map[int64]chan interface{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[int64]chan interface{}),
	}
}

func (b *Broadcaster) Subscribe(id int64) (<-chan interface{}, error) {
	defer b.mu.Unlock()
	b.mu.Lock()
	s := make(chan interface{}, 1)

	if _, ok := b.clients[id]; ok {
		return nil, fmt.Errorf("signal %d already exist", id)
	}

	b.clients[id] = s

	return b.clients[id], nil
}

func (b *Broadcaster) SubscribeWithCallback(callback func(interface{})) error {
	subscriptionId := time.Now().UnixNano()
	channel, err := b.Subscribe(subscriptionId)
	if err != nil {
		panic(err)
	}

done:
	for {
		select {
		case signal := <-channel:
			callback(signal)
		case <-context.Background().Done():
			break done
		}
	}

	b.Unsubscribe(subscriptionId)
	return nil
}

func (b *Broadcaster) Unsubscribe(id int64) {
	defer b.mu.Unlock()
	b.mu.Lock()
	if _, ok := b.clients[id]; ok {
		close(b.clients[id])
	}

	delete(b.clients, id)
}

func (b *Broadcaster) Broadcast(item interface{}) {
	defer b.mu.Unlock()
	b.mu.Lock()
	for k := range b.clients {
		if len(b.clients[k]) == 0 {
			b.clients[k] <- item
		}
	}
}
