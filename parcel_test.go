package parcel

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	// Set up test database
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcels (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			client_id INTEGER NOT NULL,
			address TEXT NOT NULL,
			status TEXT NOT NULL,
			created DATETIME NOT NULL
		)
	`)
	assert.NoError(t, err)

	store := NewParcelStore(db)

	t.Run("Register Parcel", func(t *testing.T) {
		client := 1
		address := "Test Address"

		p, err := store.Register(client, address)
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, client, p.Client)
		assert.Equal(t, address, p.Address)
		assert.Equal(t, "created", p.Status)
	})

	t.Run("Update Parcel Address", func(t *testing.T) {
		client := 2
		address := "Initial Address"
		newAddress := "Updated Address"

		p, err := store.Register(client, address)
		assert.NoError(t, err)

		err = store.UpdateAddress(p.ID, newAddress)
		assert.NoError(t, err)

		parcels, err := store.GetByClient(client)
		assert.NoError(t, err)
		assert.Len(t, parcels, 1)
		assert.Equal(t, newAddress, parcels[0].Address)
	})

	t.Run("Get Parcels by Client", func(t *testing.T) {
		client := 3
		addresses := []string{"Address 1", "Address 2", "Address 3"}

		for _, addr := range addresses {
			_, err := store.Register(client, addr)
			assert.NoError(t, err)
		}

		parcels, err := store.GetByClient(client)
		assert.NoError(t, err)
		assert.Len(t, parcels, len(addresses))

		addressMap := make(map[string]bool)
		for _, addr := range addresses {
			addressMap[addr] = true
		}

		for _, p := range parcels {
			assert.True(t, addressMap[p.Address])
		}
	})
}
