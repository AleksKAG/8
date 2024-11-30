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

func (s ParcelService) Register(client int, address string) (Parcel, error) {
	parcel := Parcel{
		Client:    client,
		Status:    ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := s.store.Add(parcel)
	if err != nil {
		return parcel, err
	}
	parcel.Number = id

	fmt.Printf("Новая посылка № %d зарегистрирована на адрес %s от клиента %d.\n", parcel.Number, address, client)
	return parcel, nil
}

func (s ParcelService) ChangeAddress(number int, address string) error {
	return s.store.SetAddress(number, address)
}

func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}

	var nextStatus string
	switch parcel.Status {
	case ParcelStatusRegistered:
		nextStatus = ParcelStatusSent
	case ParcelStatusSent:
		nextStatus = ParcelStatusDelivered
	case ParcelStatusDelivered:
		return nil
	}

	fmt.Printf("Статус посылки № %d обновлён на %s.\n", number, nextStatus)
	return s.store.SetStatus(number, nextStatus)
}

func (s ParcelService) Delete(number int) error {
	return s.store.Delete(number)
}

func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err
	}

	fmt.Printf("Посылки клиента %d:\n", client)
	for _, parcel := range parcels {
		fmt.Printf("№ %d: %s, статус: %s\n", parcel.Number, parcel.Address, parcel.Status)
	}
	return nil
}

func main() {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	service := NewParcelService(store)

	client := 1
	address := "Москва, ул. Ленина, 5"
	newAddress := "Санкт-Петербург, ул. Пушкина, 10"

	p, err := service.Register(client, address)
	if err != nil {
		log.Fatal(err)
	}

	if err := service.ChangeAddress(p.Number, newAddress); err != nil {
		log.Println(err)
	}

	if err := service.NextStatus(p.Number); err != nil {
		log.Println(err)
	}

	if err := service.PrintClientParcels(client); err != nil {
		log.Println(err)
	}

	if err := service.Delete(p.Number); err != nil {
		log.Println(err)
	}

	if err := service.PrintClientParcels(client); err != nil {
		log.Println(err)
	}
}
