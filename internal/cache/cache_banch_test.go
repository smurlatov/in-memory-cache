package cache

import (
	"github.com/google/uuid"
	"in-memory-cache/internal/models"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkCacheSet(b *testing.B) {
	cache := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := cache.Set(strconv.Itoa(i), generateProfile())
		if err != nil {
			b.Error("Failed to set value:", err)
		}
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := New()
	// Pre-fill the cache
	for i := 0; i < 100; i++ {
		cache.Set(strconv.Itoa(i), generateProfile())
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
		cache.Set(strconv.Itoa(i), generateProfile())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete(strconv.Itoa(i % 100))
	}
}

func BenchmarkCacheParallelSetGet(b *testing.B) {
	cache := New()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			profile := generateProfile()
			err := cache.Set(profile.UUID, profile)
			if err != nil {
				b.Error("Failed to set value:", err)
			}
			_, err = cache.Get(profile.UUID)
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

	time.Sleep(2 * time.Millisecond) // sleep for all values become outdated
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		cache.removeExpiredItems()
	}
}

func generateProfile() models.Profile { //generate random Profile with 0-4 orders
	ordersNumber := rand.Int() % 5
	orders := make([]*models.Order, 0, ordersNumber)
	for i := 0; i < ordersNumber; i++ {
		orders = append(orders, &models.Order{
			UUID:      uuid.New().String(),
			Value:     "order",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return models.Profile{
		UUID:   uuid.New().String(),
		Name:   "name",
		Orders: orders,
	}
}
