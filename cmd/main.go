package main

import (
	"github.com/google/uuid"
	"in-memory-cache/internal/cache"
	"in-memory-cache/internal/models"
	"math/rand"
	"time"
)

func main() {
	inMemCache := cache.New()
	profileUUIDList := make([]string, 0, 10)
	for i := 0; i < 10; i++ { // generate profiles map , that will use as a fixture
		profileUUIDList = append(profileUUIDList, uuid.New().String())
	}

	for i := 0; i < 40; i++ {
		go func() {
			inMemCache.Set(profileUUIDList[rand.Int()%10], generateProfile()) //get random Id from list and Set value by this id
			inMemCache.Get(profileUUIDList[rand.Int()%10])
		}()
	}
	time.Sleep(time.Second * 15)
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
