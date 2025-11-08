package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare

	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add

	id, err := store.Add(parcel)

	require.NoError(t, err)

	assert.NotEqual(t, 0, id, "The parcel has not been added")

	// get

	p, err := store.Get(id)

	require.NoError(t, err)

	assert.Equal(t, parcel.Client, p.Client, "The Client does not match")
	assert.Equal(t, parcel.Status, p.Status, "The Status does not match")
	assert.Equal(t, parcel.Address, p.Address, "The Address does not match")
	assert.Equal(t, parcel.CreatedAt, p.CreatedAt, "The CreatedAt does not match")

	// delete

	err = store.Delete(id)
	require.NoError(t, err)

	pp, err := store.Get(id)
	require.NoError(t, err)

	assert.NotEqual(t, p, pp, "The parcel has not been removed from the database")

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add

	id, err := store.Add(parcel)

	require.NoError(t, err)

	assert.NotEqual(t, 0, id, "The parcel has not been added")

	// set address

	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)

	require.NoError(t, err)

	// check

	p, err := store.Get(id)

	require.NoError(t, err)

	assert.Equal(t, newAddress, p.Address, "The parcel Address has not been updated")

}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add

	id, err := store.Add(parcel)

	require.NoError(t, err)

	assert.NotEqual(t, 0, id, "The parcel has not been added")

	// set status

	newStatus := ParcelStatusSent

	err = store.SetStatus(id, newStatus)

	require.NoError(t, err)

	// check

	p, err := store.Get(id)

	require.NoError(t, err)

	assert.Equal(t, newStatus, p.Status, "The parcel Status has not been updated")

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {

		id, err := store.Add(parcels[i])

		require.NoError(t, err)

		assert.NotEqual(t, 0, id, "The parcel has not been added")

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client

	// убедитесь в отсутствии ошибки
	require.NoError(t, err)

	// проверяем, что количество полученных посылок совпадает с количеством добавленных
	assert.Equal(t, len(parcels), len(storedParcels), "The number of parcels does not match")

	// check
	for _, parcel := range storedParcels {

		p := parcelMap[parcel.Number]
		assert.Equal(t, p.Number, parcel.Number, "The Number does not match")

		assert.Equal(t, p.Client, parcel.Client, "The Client does not match")
		assert.Equal(t, p.Status, parcel.Status, "The Status does not match")
		assert.Equal(t, p.Address, parcel.Address, "The Address does not match")
		assert.Equal(t, p.CreatedAt, parcel.CreatedAt, "The CreatedAt does not match")

	}
}
