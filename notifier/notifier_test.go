package notifier

import (
	"sync"
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {
	const numSubscribers = 10
	n := NewNotifier()
	defer cleanup(n)
	wg := &sync.WaitGroup{}

	done := make(chan bool)
	wg.Add(numSubscribers)
	go func() {
		wg.Wait()
		close(done)
	}()

	for i := 0; i < numSubscribers; i++ {
		go func(wg *sync.WaitGroup) {
			ch := make(chan []byte)
			n.Subscribe(ch)
			defer wg.Done()
		}(wg)
	}

	<-done

	if len(n.clientMap) != numSubscribers {
		t.Errorf("Number of subscribers should be %d but it's %d", numSubscribers, len(n.clientMap))
	}
}

func TestUnsubscribe(t *testing.T) {
	n := NewNotifier()
	ch := make(chan []byte)

	id := n.Subscribe(ch)

	n.Unsubscribe(id)

	select {
	case <-ch:
	default:
		t.Errorf("channel should be closed")
	}

	if len(n.clientMap) > 0 {
		t.Errorf("clientMap should be empty")
	}
}

func TestNotifyAll(t *testing.T) {
	const numSubscribers = 10
	n := NewNotifier()
	defer cleanup(n)

	channels := make([]chan []byte, numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		ch := make(chan []byte, 1)
		channels[i] = ch
		n.Subscribe(ch)
	}

	if len(n.clientMap) != numSubscribers {
		t.Errorf("Number of subscribers should be %d but it's %d", numSubscribers, len(n.clientMap))
	}

	n.NotifyAll([]byte("message"))
	time.Sleep(time.Millisecond * 10)

	for i := 0; i < numSubscribers; i++ {
		select {
		case msg := <-channels[i]:
			if string(msg) != "message" {
				t.Errorf(`msg from channel #%d should be "message", but it's "%s"`, i, string(msg))
			}
		default:
			t.Errorf("client %d Should received a message", i)

		}
	}
}

func TestNotifyOne(t *testing.T) {
	n := NewNotifier()
	defer cleanup(n)

	ch := make(chan []byte, 1)
	id := n.Subscribe(ch)

	n.NotifyOne(id, []byte("message"))
	wait := time.After(time.Millisecond * 200)

	select {
	case msg := <-ch:
		if string(msg) != "message" {
			t.Errorf(`msg from channel should be "message", but it's "%s"`, string(msg))
		}

	case <-wait:
		t.Error("client Should received a message")
	}
}

func cleanup(n *Notifier) {
	for id := range n.clientMap {
		n.Unsubscribe(id)
	}
}
