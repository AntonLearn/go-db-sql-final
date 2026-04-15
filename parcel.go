package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	result, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) "+
		"VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client), sql.Named("status", p.Status),
		sql.Named("address", p.Address), sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	idLast, err := result.LastInsertId()
	return int(idLast), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	row := s.db.QueryRow("SELECT number, client, status, address, created_at "+
		"FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel
	parcels, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel "+
		"WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return res, err
	}
	defer parcels.Close()
	for parcels.Next() {
		p := Parcel{}
		err := parcels.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}
	err = parcels.Err()
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, nextStatus string) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}
	if p.Status == ParcelStatusDelivered {
		return ErrStatusWhenDelivered
	}
	_, err = s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", nextStatus), sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}
	if p.Status != ParcelStatusRegistered {
		return ErrResetWhenNotRegistered
	}
	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND "+
		"status = :status", sql.Named("address", address), sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}
	if p.Status != ParcelStatusRegistered {
		return ErrDeleteWhenNotRegistered
	}
	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}
