package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	query := `
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		);
	`
	_, err = db.Exec(query)
	require.NoError(t, err)

	return db
}

func TestParcelService(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)
	service := NewParcelService(store)

	client := 1
	address := "Тестовый адрес"
	newAddress := "Обновленный адрес"

	// Test Register
	p, err := service.Register(client, address)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusRegistered, p.Status)

	// Test Change Address
	err = service.ChangeAddress(p.Number, newAddress)
	require.NoError(t, err)

	// Test GetByClient
	parcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, 1, len(parcels))
	require.Equal(t, newAddress, parcels[0].Address)

	// Test Delete
	err = service.Delete(p.Number)
	require.NoError(t, err)

	parcels, err = store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, 0, len(parcels))
}
