package cache

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkCacheSet(b *testing.B) {
	cache := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := cache.Set(strconv.Itoa(i), map[string]int{"number": i})
		if err != nil {
			b.Error("Failed to set value:", err)
		}
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := New()
	// Pre-fill the cache
	for i := 0; i < 100; i++ {
		cache.Set(strconv.Itoa(i), map[string]int{"number": i})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cache.Get(strconv.Itoa(i % 100)) // Ensure keys exist
		if err != nil {
			b.Error("Failed to get value:", err)
		}
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	cache := New()
	// Pre-fill the cache
	for i := 0; i < 100; i++ {
		cache.Set(strconv.Itoa(i), map[string]int{"number": i})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete(strconv.Itoa(i % 100)) // Repeatedly delete and re-add
		cache.Set(strconv.Itoa(i%100), map[string]int{"number": i})
	}
}

func BenchmarkCacheParallelSetGet(b *testing.B) {
	cache := New()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			key := strconv.Itoa(i)
			err := cache.Set(key, map[string]int{"number": i})
			if err != nil {
				b.Error("Failed to set value:", err)
			}
			_, err = cache.Get(key)
			if err != nil {
				b.Error("Failed to get value:", err)
			}
		}
	})
}

func BenchmarkCacheGC(b *testing.B) {
	cache := New(1 * time.Millisecond) //set small ttl for test

	for i := 0; i < 10000; i++ {
		cache.Set(strconv.Itoa(i), map[string]int{"number": i})
	}

	b.ResetTimer()
	time.Sleep(2 * time.Millisecond)
	for n := 0; n < b.N; n++ {
		cache.removeExpiredItems()
	}
}
