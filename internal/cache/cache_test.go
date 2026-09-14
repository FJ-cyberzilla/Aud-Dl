package cache

import (
	"testing"
	"time"
)

func TestCache_SetGet(t *testing.T) {
	c := NewCacheManager(1 * time.Second)
	key := "test"
	value := "value"
	
	c.Set(key, value, 1*time.Second)
	
	val, found := c.Get(key)
	if !found {
		t.Errorf("Expected to find key %s", key)
	}
	if val != value {
		t.Errorf("Expected value %v, got %v", value, val)
	}
}

func TestCache_Expired(t *testing.T) {
	c := NewCacheManager(1 * time.Second)
	key := "expired"
	
	c.Set(key, "data", 10*time.Millisecond)
	
	time.Sleep(20 * time.Millisecond)
	
	_, found := c.Get(key)
	if found {
		t.Errorf("Expected key %s to be expired", key)
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCacheManager(1 * time.Second)
	key := "delete"
	
	c.Set(key, "data", 0)
	c.Delete(key)
	
	_, found := c.Get(key)
	if found {
		t.Errorf("Expected key %s to be deleted", key)
	}
}
