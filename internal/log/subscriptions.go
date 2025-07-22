package log

import "sync"

type Subscription struct {
	K int
	C chan Line
}

type Subscribers struct {
	mu          sync.Mutex
	nextID      int
	Subscribers []*Subscription
}

// Subscribe creates a new subscription in a concurrency-safe way
func (s *Subscribers) Subscribe() *Subscription {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub := &Subscription{
		K: s.nextID,
		C: make(chan Line, 100), // buffered to avoid blocking
	}

	s.nextID++
	s.Subscribers = append(s.Subscribers, sub)
	return sub
}

// Unsubscribe removes a given subscription from the list and closes its channel
func (s *Subscribers) Unsubscribe(target *Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var updated []*Subscription
	for _, ss := range s.Subscribers {
		if ss != target {
			updated = append(updated, ss)
		}
	}
	close(target.C)
	s.Subscribers = updated
}

// Filter returns a filtered slice of subscriptions matching a predicate
func (s *Subscribers) Filter(predicate func(k int, ss *Subscription) bool) []*Subscription {
	s.mu.Lock()
	defer s.mu.Unlock()

	var filtered []*Subscription
	for _, ss := range s.Subscribers {
		if predicate(ss.K, ss) {
			filtered = append(filtered, ss)
		}
	}
	return filtered
}

// Broadcast sends a log line to all subscribers
func (s *Subscribers) Broadcast(line Line) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ss := range s.Subscribers {
		if ss.C != nil {
			ss.C <- line
		}
	}
}
