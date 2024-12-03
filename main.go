package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

const (
	ParcelStatusRegistered = "registered"
	ParcelStatusSent       = "sent"
	ParcelStatusDelivered  = "delivered"
)

type Parcel struct {
	Number    int
	Client    int
	Status    string
	Address   string
	CreatedAt string
}

type ParcelService struct {
	store ParcelStore
}

func NewParcelService(store ParcelStore) ParcelService {
	return ParcelService{store: store}
}

// Остальные методы ParcelService остаются без изменений...

func main() {
	// Настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db) // создайте объект ParcelStore функцией NewParcelStore
	service := NewParcelService(store)

	// Пример использования
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parcel created: %+v\n", p)
}
