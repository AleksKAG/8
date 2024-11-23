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

func main() {
	// Подключение к базе данных
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Создание таблицы, если она ещё не создана
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		)
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	store := NewParcelStore(db)
	service := NewParcelService(store)

	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(p.Number, newAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.NextStatus(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	p, err = service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
