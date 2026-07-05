package sse

import "sync"

type Hub struct {
	mu      sync.Mutex
	clients map[chan string]struct{}
}

var defaultHub = &Hub{
	clients: make(map[chan string]struct{}),
}

func Default() *Hub { return defaultHub }

func (h *Hub) Subscribe() chan string {
	ch := make(chan string, 4)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch chan string) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

func (h *Hub) Broadcast(event string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- event:
		default:
		}
	}
}
