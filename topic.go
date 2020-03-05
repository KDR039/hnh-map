package main

import "sync"

type topic struct {
	c  []chan *TileData
	mu sync.Mutex
}

func (t *topic) watch(c chan *TileData) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.c = append(t.c, c)
}

func (t *topic) send(b *TileData) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, c := range t.c {
		select {
		case c <- b:
		default:
			close(c)
			t.c[i] = t.c[len(t.c)-1]
			t.c = t.c[:len(t.c)-1]
		}
	}
}

func (t *topic) close() {
	for _, c := range t.c {
		close(c)
	}
	t.c = t.c[:0]
}