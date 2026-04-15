package main

import (
	"database/sql"
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
	log.Printf("New parcel № %d on address %s from client with id %d "+
		"registered at %s\n", parcel.Number, parcel.Address, parcel.Client,
		parcel.CreatedAt)
	return parcel, nil
}

func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err
	}
	log.Printf("Parcels of client id %d:\n", client)
	for _, parcel := range parcels {
		log.Printf("Parcel № %d on address %s from client id %d registered at %s in status %s\n",
			parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt, parcel.Status)
	}
	return nil
}

func (s ParcelService) NextStatus(number int) error {
	err := s.store.SetStatus(number)
	if err != nil {
		return err
	}
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}
	log.Printf("Parcel № %d on address %s from client with id %d registered at %s: "+
		"has new status: %s\n", parcel.Number, parcel.Address, parcel.Client,
		parcel.CreatedAt, parcel.Status)
	return nil
}

func (s ParcelService) ChangeAddress(number int, address string) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}
	oldAddress := parcel.Address
	err = s.store.SetAddress(number, address)
	if err != nil {
		return err
	}
	parcel, err = s.store.Get(number)
	if err != nil {
		return err
	}
	log.Printf("Parcel № %d on address %s from client with id %d registered at %s: "+
		"change delivery address. New delivery address is %s\n",
		parcel.Number, oldAddress, parcel.Client, parcel.CreatedAt, parcel.Address)
	return nil
}

func (s ParcelService) Delete(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}
	err = s.store.Delete(number)
	if err != nil {
		return err
	}
	log.Printf("Parcel № %d on address %s from client with id %d registered at %s has been removed\n",
		parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt)
	return nil
}

func main() {
	// connection to db
	dbName := "tracker.db"
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	// creating object ParcelStore by function NewParcelStore
	store := NewParcelStore(db)
	// creating object ParcelService by function NewParcelService
	service := NewParcelService(store)
	// first parcel registration
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		log.Println(err)
	}
	number := p.Number
	// change first parcel address
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(number, newAddress)
	if err != nil {
		log.Println(err)
	}
	// setting next first parcel status
	err = service.NextStatus(number)
	if err != nil {
		log.Println(err)
	}
	// print list of client parcels
	err = service.PrintClientParcels(client)
	if err != nil {
		log.Println(err)
	}
	// removing first parcel
	err = service.Delete(number)
	if err != nil {
		log.Println(err)
	}
	// print list of client parcels
	// first parcel was not removed because its status was not registered
	err = service.PrintClientParcels(client)
	if err != nil {
		log.Println(err)
	}
	// second parcel registration
	p, err = service.Register(client, address)
	if err != nil {
		log.Println(err)
	}
	number = p.Number
	// removing second parcel
	err = service.Delete(number)
	if err != nil {
		log.Println(err)
	}
	// print list of client parcels
	// second parcel was removed because its status was registered
	err = service.PrintClientParcels(client)
	if err != nil {
		log.Println(err)
	}
}
