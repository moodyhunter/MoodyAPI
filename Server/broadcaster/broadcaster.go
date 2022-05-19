package broadcaster

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Broadcaster[T any] struct {
	lock    sync.Mutex
	clients map[int64]chan *T
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		clients: make(map[int64]chan *T),
	}
}

func (b *Broadcaster[T]) Subscribe(id int64) (<-chan *T, error) {
	defer b.lock.Unlock()
	b.lock.Lock()
	s := make(chan *T, 1)

	if _, ok := b.clients[id]; ok {
		return nil, fmt.Errorf("signal %d already exist", id)
	}

	b.clients[id] = s
	return b.clients[id], nil
}

func (b *Broadcaster[T]) BlockedSubscribeWithCallback(callback func(*T)) error {
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

func (b *Broadcaster[T]) Unsubscribe(id int64) {
	defer b.lock.Unlock()
	b.lock.Lock()
	if _, ok := b.clients[id]; ok {
		close(b.clients[id])
	}

	delete(b.clients, id)
}

func (b *Broadcaster[T]) Broadcast(item *T) {
	defer b.lock.Unlock()
	b.lock.Lock()
	for k := range b.clients {
		if len(b.clients[k]) == 0 {
			b.clients[k] <- item
		}
	}
}
