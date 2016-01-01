package cache

import (
	"testing"
	"container/heap"
	"time"
	"fmt"
	//"github.com/bradfitz/gomemcache/memcache"
)

func TestHeap(t *testing.T) {
	//Testing the priority queue on expirations
	var h expirationHeap
	h = make([]*expirationItem, 0)
	heap.Init(&h)
	
	heap.Push(&h, &expirationItem{key:"1", expiration:time.Unix(100, 0)})
	heap.Push(&h, &expirationItem{key:"2", expiration:time.Unix(20, 0)})
	heap.Push(&h, &expirationItem{key:"3", expiration:time.Unix(50, 0)})
	heap.Push(&h, &expirationItem{key:"4", expiration:time.Unix(500, 0)})
	
	
	for i, item := range h {
		if item.index != i {
			t.Errorf("index mismatch: wanted %d, got %d", i, item.index)
		}
	}
	
	if h[0].key != "2" {
		t.Errorf("min in heap should be %s, got %s", "2", h[0].key)
	}
	
	item := heap.Pop(&h).(*expirationItem)
	if item.key != "2" {
		t.Errorf("min poped from heap should be %s, got %s", "2", item.key)
	}
	
	for i, item := range h {
		if item.index != i {
			t.Errorf("index mismatch: wanted %d, got %d", i, item.index)
		}
	}
	
	if h[0].key != "3" {
		t.Errorf("min in heap should be %s, got %s", "3", h[0].key)
	}
	
	for _, item := range h {
		if item.key == "4" {
			item.expiration = time.Unix(1, 0)
			heap.Fix(&h, item.index)
		}
	}
	
	for i, item := range h {
		if item.index != i {
			t.Errorf("index mismatch: wanted %d, got %d", i, item.index)
		}
	}
	
	if h[0].key != "4" {
		t.Errorf("min in heap should be %s, got %s", "4", h[0].key)
	}
	
	
	item = heap.Pop(&h).(*expirationItem)
	if item.key != "4" {
		t.Errorf("poped value in heap should be %s, got %s", "4", item.key)
	}
	
	for i, item := range h {
		if item.index != i {
			t.Errorf("index mismatch: wanted %d, got %d", i, item.index)
		}
	}
	
	item = heap.Pop(&h).(*expirationItem)
	if item.key != "3" {
		t.Errorf("poped value in heap should be %s, got %s", "3", item.key)
	}
	
	
	for i, item := range h {
		if item.index != i {
			t.Errorf("index mismatch: wanted %d, got %d", i, item.index)
		}
	}
	
	item = heap.Pop(&h).(*expirationItem)
	if item.key != "1" {
		t.Errorf("poped value in heap should be %s, got %s", "1", item.key)
	}
	
	for i, item := range h {
		if item.index != i {
			t.Errorf("index mismatch: wanted %d, got %d", i, item.index)
		}
	}
	
	if h.Len() != 0 {
		t.Error("Should be empty heap now")
	}
}

func TestSimplecache(t *testing.T) {
	s := NewSimplecache()
	
	
	// try insert an expired item
	err := s.Set("1", []byte("1"), -1)
	if fmt.Sprintf("%s", err) != "Expired" {
		t.Errorf("Get should return expired error, but got %s", err)
	}
	
	_, err = s.Get("1")
	if fmt.Sprintf("%s", err) != "Not found" {
		t.Errorf("Get should return not found error, but got %s", err)
	}
	
	// insert item, that never expires
	err = s.Set("0", []byte("0"), 0)
	if err != nil {
		t.Errorf("Cache insertion returns error: %s", err)
	}
	// insert another item (1 second from now to expire)
	err = s.Set("1", []byte("1"), 1)
	if err != nil {
		t.Errorf("Cache insertion returns error: %s", err)
	}
	// insert 2nd item (2 seconds from now to expire)
	err = s.Set("2", []byte("2"), 2)
	if err != nil {
		t.Errorf("Cache insertion returns error: %s", err)
	}
	// insert 3rd item (3 seconds from now to expire)
	err = s.Set("3", []byte("3"), 3)
	if err != nil {
		t.Errorf("Cache insertion returns error: %s", err)
	}
	
	// get these items, should all return no error
	item, err := s.Get("1")
	if err != nil {
		t.Errorf("Get returns error: %s", err)
	}
	if string(item) != "1" {
		t.Errorf("Get returns value should be 1, got %s", string(item))
	}
	
	item, err = s.Get("2")
	if err != nil {
		t.Errorf("Get returns error: %s", err)
	}
	if string(item) != "2" {
		t.Errorf("Get returns value should be 2, got %s", string(item))
	}
	
	item, err = s.Get("3")
	if err != nil {
		t.Errorf("Get returns error: %s", err)
	}
	if string(item) != "3" {
		t.Errorf("Get returns value should be 3, got %s", string(item))
	}
	
	// after one second, the first should be purged
	time.Sleep(time.Second)
	
	_, err = s.Get("1")
	if fmt.Sprintf("%s", err) != "Expired" {
		t.Errorf("Get should return expired error, but got %s", err)
	}
	
	time.Sleep(time.Duration(100) * time.Millisecond)
	_, err = s.Get("1")
	if fmt.Sprintf("%s", err) != "Not found" {
		t.Errorf("Get should return not found error, but got %s", err)
	}
	
	// note that not only cache should be purged, the expiration heap
	if s.expirations.Len() != 3 {
		t.Errorf("The expiration heap not purged correctly")
	}
	
	item, err = s.Get("2")
	if err != nil {
		t.Errorf("Get returns error: %s", err)
	}
	if string(item) != "2" {
		t.Errorf("Get returns value should be 2, got %s", string(item))
	}
	
	item, err = s.Get("3")
	if err != nil {
		t.Errorf("Get returns error: %s", err)
	}
	if string(item) != "3" {
		t.Errorf("Get returns value should be 3, got %s", string(item))
	}
	
	//try delete
	s.Delete("3")
	_, err = s.Get("3")
	if fmt.Sprintf("%s", err) != "Not found" {
		t.Errorf("Get should return not found error, but got %s", err)
	}
	
	if s.expirations.Len() != 2 {
		t.Errorf("The expiration heap not purged correctly")
	}
	
	_, err = s.Get("0")
	if err != nil {
		t.Errorf("Get returns error: %s", err)
	}
}

// comment this out
// Because we do not always have local memcache server setup
/*
func TestMemcache(t *testing.T) {
	s := NewMemcache("127.0.0.1:11211")
	
	// try insert an expired item
	err := s.Set("1", []byte("1"), -1)	
	_, err = s.Get("1")
	if err != memcache.ErrCacheMiss {
		t.Errorf("Get should return %v, but got %v", memcache.ErrCacheMiss, err)
	}
	
	// insert item, that never expires
	err = s.Set("0", []byte("0"), 0)
	if err != nil {
		t.Errorf("Cache insertion returns error: %v", err)
	}
	// insert another item (1 second from now to expire)
	err = s.Set("1", []byte("1"), 1)
	if err != nil {
		t.Errorf("Cache insertion returns error: %v", err)
	}
	// insert 2nd item (2 seconds from now to expire)
	err = s.Set("2", []byte("2"), 2)
	if err != nil {
		t.Errorf("Cache insertion returns error: %v", err)
	}
	// insert 3rd item (3 seconds from now to expire)
	err = s.Set("3", []byte("3"), 3)
	if err != nil {
		t.Errorf("Cache insertion returns error: %v", err)
	}
	
	// get these items, should all return no error
	item, err := s.Get("1")
	if err != nil {
		t.Errorf("Get returns error: %v", err)
	}
	if string(item) != "1" {
		t.Errorf("Get returns value should be 1, got %s", string(item))
	}
	
	item, err = s.Get("2")
	if err != nil {
		t.Errorf("Get returns error: %v", err)
	}
	if string(item) != "2" {
		t.Errorf("Get returns value should be 2, got %s", string(item))
	}
	
	item, err = s.Get("3")
	if err != nil {
		t.Errorf("Get returns error: %v", err)
	}
	if string(item) != "3" {
		t.Errorf("Get returns value should be 3, got %s", string(item))
	}
	
	// after one second, the first should be purged
	time.Sleep(time.Second)
	
	_, err = s.Get("1")
	if err != memcache.ErrCacheMiss {
		t.Errorf("Get should return %v, but got %v", memcache.ErrCacheMiss, err)
	}
	
	item, err = s.Get("2")
	if err != nil {
		t.Errorf("Get returns error: %v", err)
	}
	if string(item) != "2" {
		t.Errorf("Get returns value should be 2, got %s", string(item))
	}
	
	item, err = s.Get("3")
	if err != nil {
		t.Errorf("Get returns error: %v", err)
	}
	if string(item) != "3" {
		t.Errorf("Get returns value should be 3, got %s", string(item))
	}
	
	//try delete
	s.Delete("3")
	_, err = s.Get("3")
	if err != memcache.ErrCacheMiss {
		t.Errorf("Get should return %v, but got %v", memcache.ErrCacheMiss, err)
	}
		
	_, err = s.Get("0")
	if err != nil {
		t.Errorf("Get returns error: %v", err)
	}
}*/