package main

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"in-memory-cache/internal/cache"
	"in-memory-cache/internal/models"
	"math/rand"
	"time"
)

func main() {
	startTime := time.Now() // Запоминаем время начала выполнения
	cache := cache.New()
	newUuid := uuid.New().String()
	err := cache.Set(newUuid, generateProfile())
	fmt.Println(err)
	fmt.Println(cache.Get(newUuid))
	//time.Sleep(1 * time.Second)
	//fmt.Println(cache.Get(newUuid))
	//time.Sleep(2 * time.Second)
	//fmt.Println(cache.Get(newUuid))
	//time.Sleep(10 * time.Second)
	elapsedTime := time.Since(startTime) // Вычисляем затраченное время
	fmt.Printf("Программа выполнялась: %s\n", elapsedTime)
}

func generateProfile() models.Profile { //generate random Profile with 0-4 orders
	ordersNumber := rand.Int() % 5
	fmt.Println(ordersNumber)
	orders := make([]*models.Order, 0, ordersNumber)
	for i := 0; i < ordersNumber; i++ {
		orders = append(orders, &models.Order{
			UUID:      uuid.New().String(),
			Value:     faker.Name(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return models.Profile{
		UUID:   uuid.New().String(),
		Name:   faker.Name(),
		Orders: orders,
	}
}
