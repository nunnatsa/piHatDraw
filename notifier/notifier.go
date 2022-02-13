package notifier

import (
	"log"
	"sync"
	"sync/atomic"
)

type idProvider struct {
	counter uint64
}

func newIdProvider() *idProvider {
	return &idProvider{
		counter: 0,
	}
}

func (p *idProvider) getNextID() uint64 {
	return atomic.AddUint64(&p.counter, 1)
}

type Notifier struct {
	clientMap map[uint64]chan []byte
	idp       *idProvider
	lock      *sync.Mutex
}

func NewNotifier() *Notifier {
	return &Notifier{
		clientMap: make(map[uint64]chan []byte),
		idp:       newIdProvider(),
		lock:      &sync.Mutex{},
	}
}

func (n *Notifier) Subscribe(ch chan []byte) uint64 {
	id := n.idp.getNextID()
	n.lock.Lock()
	n.clientMap[id] = ch
	n.lock.Unlock()

	log.Println("register new client", id)

	return id
}

func (n *Notifier) Unsubscribe(id uint64) {
	log.Println("deregister client", id)
	if ch, ok := n.clientMap[id]; ok {
		n.lock.Lock()
		delete(n.clientMap, id)
		n.lock.Unlock()
		close(ch)
	}
}

func (n *Notifier) NotifyAll(data []byte) {
	for _, subscriber := range n.clientMap {
		n.sendToSubscriber(subscriber, data)
	}
}

func (n *Notifier) NotifyOne(id uint64, data []byte) {
	if subscriber, ok := n.clientMap[id]; ok {
		n.sendToSubscriber(subscriber, data)
	} else {
		log.Printf("subscriber id %d was not found", id)
	}
}

func (n Notifier) sendToSubscriber(subscriber chan<- []byte, data []byte) {
	subscriber <- data
}

func (n Notifier) Close() {
	for id := range n.clientMap {
		n.Unsubscribe(id)
	}
}