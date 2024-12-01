package parcel

import (
	"database/sql"
	"fmt"
	"time"
)

type Parcel struct {
	ID      int
	Client  int
	Address string
	Status  string
	Created time.Time
}

type Storage struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Register(client int, address string) (*Parcel, error) {
	query := `INSERT INTO parcels (client_id, address, status, created) VALUES (?, ?, 'created', ?) RETURNING id, created`
	row := s.db.QueryRow(query, client, address, time.Now())

	var p Parcel
	p.Client = client
	p.Address = address
	p.Status = "created"
	if err := row.Scan(&p.ID, &p.Created); err != nil {
		return nil, fmt.Errorf("failed to register parcel: %w", err)
	}

	return &p, nil
}

func (s *Storage) UpdateAddress(parcelID int, newAddress string) error {
	query := `UPDATE parcels SET address = ? WHERE id = ?`
	result, err := s.db.Exec(query, newAddress, parcelID)
	if err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no parcel found with id %d", parcelID)
	}
	return nil
}

func (s *Storage) GetByClient(clientID int) ([]Parcel, error) {
	query := `SELECT id, address, status, created FROM parcels WHERE client_id = ?`
	rows, err := s.db.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.ID, &p.Address, &p.Status, &p.Created); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		parcels = append(parcels, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return parcels, nil
}
