package cache

import (
	"sync"
	"time"
	"container/heap"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	MaxUnixTime int64 = 9999999999
)
// We define a cache interface same as memcached golang client
// By this way, we can interchange local caches & memcached with
// interface unchanged.

// We only provide interfaces that might be relavent to token cache

type Cache interface {
	// Set with key, value, and expiration in seconds
	Set(string, []byte, int) error
	// Get the value given the key
	Get(string) ([]byte, error)
	// Delete the value stored in the key
	Delete(string) error
}

// Used to save priority queue-based structure for purging
// expired keys
// key is required because we need a way to find the store[key]
// from expirationItem
type expirationItem struct {
	key string
	expiration time.Time
	index int
}
type expirationHeap []*expirationItem
func (h expirationHeap) Len() int           { return len(h) }
func (h expirationHeap) Less(i, j int) bool { return h[i].expiration.Before(h[j].expiration) }
func (h expirationHeap) Swap(i, j int)      { 
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *expirationHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*expirationItem)
	item.index = n
	*h = append(*h, item)
}
func (h *expirationHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1 //for safety
	*h = old[0 : n-1]
	return item
}

// A simple local cache to store tokens if memcache servers not present.
// One design decision is how to check which key expires and needs to be purged.
// One way is to use another thread and check every some fixed interval.
// Another way is to call purge after Set or Get is called (one more is saved in.) We
// use the latter one for simplicity. The purge will be called whenever a Set is 
// called, or whenenver a Get is called but the got item is expired.
type simplecacheItem struct {
	item []byte
	//We need a way to find the expiration item on priority-queue
	//from the simpleCacheItem.
	expiration *expirationItem
}
type Simplecache struct {
	// purge expired token shall be performed in another go thread
	// so use mutex to sync
	sync.Mutex
	
	// data store
	store map[string]*simplecacheItem 
	// used to check which key expired
	expirations expirationHeap
	// if true, another threading is doing purge
	// we shall ensure that only one thread is doing purge at
	// any time
	purging bool
}
func NewSimplecache() *Simplecache {
	cache := &Simplecache {
		store: make(map[string]*simplecacheItem),
		expirations: expirationHeap(make([]*expirationItem, 0)),
	}
	heap.Init(&cache.expirations)
	return cache
}
func (s *Simplecache) purge() {
	s.Lock()
	if s.purging {
		s.Unlock()
		return
	}
	s.Unlock()
	
	// start purging, first retrieve expired keyitems, then purge
	// these keyitems from store, use a go thread to do that
	go func() {
		// during purging, the structure has to be locked for some time
		// this may reduce system performance
		s.Lock()
		defer s.Unlock()
		
		s.purging = true
		for len(s.expirations) != 0 {
			minItem := s.expirations[0]
			if minItem.expiration.After(time.Now()) {
				break;
			}
			_ = s.expirations.Pop()
			
			//move the item from map too
			delete(s.store, minItem.key)
		}
		s.purging = false
	}()
}
func (s *Simplecache) Set(key string, value []byte, e int) error {
	//firstly, try to see whether already existed
	var expiration time.Time
	if e < 0 {
		return errors.New("Expired")
	} else if e == 0 {
		//max value, never expires
		expiration = time.Unix(MaxUnixTime,0)
	} else {
		expiration = time.Now().Add(time.Duration(e) * time.Second)
	}
	
	s.Lock()
	if _, ok := s.store[key]; ok {
		//already existed
		//change the value, and also update the expiration item
		s.store[key].item = value
		s.store[key].expiration.expiration = expiration
		heap.Fix(&s.expirations, s.store[key].expiration.index)
	} else {
		item := &expirationItem{ key: key, expiration: expiration }
		heap.Push(&s.expirations, item)
		s.store[key] = &simplecacheItem{ item: value, expiration: item }
	}
	s.Unlock()
	
	s.purge()
	
	return nil
}
func (s *Simplecache) Get(key string) ([]byte, error) {
	s.Lock()
	if s.store != nil {
		if value, ok := s.store[key]; ok {
			item := value.expiration
			if item.expiration.Before(time.Now()) {
				//expired
				heap.Remove(&s.expirations, item.index)
				delete(s.store, key)
				
				s.Unlock()
				s.purge()
				return nil, errors.New("Expired")
			}
			s.Unlock()
			return value.item, nil
		}
	}
	s.Unlock()
	return nil, errors.New("Not found")
	
}
func (s *Simplecache) Delete(key string) (error) {
	s.Lock()
	defer s.Unlock()
	
	if value, ok := s.store[key]; ok {
		heap.Remove(&s.expirations, value.expiration.index)
		delete(s.store, key)
	}
	
	return nil
}

// A wrapper to memcache client
type Memcache struct {
	client *memcache.Client
}

func NewMemcache(server ...string) *Memcache {
	return &Memcache{
		client: memcache.New(server ...),
	}
}


// Memcached client:
//  Item is an item to be got or stored in a memcached server.
//	type Item struct {
//		// Key is the Item's key (250 bytes maximum).
//		Key string
//	
//		// Value is the Item's value.
//		Value []byte
//	
//		// Object is the Item's value for use with a Codec.
//		Object interface{}
//	
//		// Flags are server-opaque flags whose semantics are entirely
//		// up to the app.
//		Flags uint32
//	
//		// Expiration is the cache expiration time, in seconds: either a relative
//		// time from now (up to 1 month), or an absolute Unix epoch time.
//		// Zero means the Item has no expiration time.
//		Expiration int32
//	
//		// Compare and swap ID.
//		casid uint64
//	}
func (m *Memcache) Set(key string, value []byte, e int) error {
	err := m.client.Set(&memcache.Item{
		Key: key,
		Value: value,
		Expiration: int32(e),
	})
	return err
}


func (m *Memcache) Get(key string) ([]byte, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return nil, err
	} else {
		return item.Value, nil
	}
}

func (m *Memcache) Delete(key string) error {
	err := m.client.Delete(key)
	return err
}