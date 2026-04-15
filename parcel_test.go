package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	// randSource is source of pseudo-random numbers
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange use randSource for generating pseudo-random numbers
	randRange = rand.New(randSource)
)

// getTestParcel generate test parcel
func getTestParcel() Parcel {
	return Parcel{
		Client:    randRange.Intn(9_999_999) + 1, // generate pseudo-random number in range from 1 to 10_000_000
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete tests add, get and remove of test parcel
func TestAddGetDelete(t *testing.T) {
	// connection to db
	dbName := "tracker.db"
	db, err := sql.Open("sqlite", dbName)
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	testParcel := getTestParcel() // generate test parcel
	// add new parcel, testing for error and not empty id
	id, err := store.Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	testParcel.Number = id
	// get this parcel by id, testing for error and it value equal test parcel value
	parcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, testParcel, parcel)
	// delete this parcel by id, testing for error and that parcel is not found in db
	err = store.Delete(id)
	require.NoError(t, err)
	// attempt getting this parcel by id, testing that will return error
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress tests on updating address test parcel
func TestSetAddress(t *testing.T) {
	// connection to db
	dbName := "tracker.db"
	db, err := sql.Open("sqlite", dbName)
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	testParcel := getTestParcel() // generate test parcel
	// add new parcel, testing for error and not empty id
	id, err := store.Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// update address, testing for error
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)
	// tests updating of address test parcel
	// get this parcel by id, testing for error and that getting address of parcel equal new address test parcel
	parcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, parcel.Address)
}

// TestSetStatus tests on updating parcel status
func TestSetStatus(t *testing.T) {
	// connection to db
	dbName := "tracker.db"
	db, err := sql.Open("sqlite", dbName)
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	testParcel := getTestParcel() // generate test parcel
	// add new parcel, testing for error and not empty id
	id, err := store.Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// update status, testing for error
	var nextStatus string
	switch testParcel.Status {
	case ParcelStatusRegistered:
		nextStatus = ParcelStatusSent
	case ParcelStatusSent:
		nextStatus = ParcelStatusDelivered
	default:
		nextStatus = ParcelStatusDelivered
	}
	err = store.SetStatus(id, nextStatus)
	require.NoError(t, err)
	// tests updating of status test parcel
	// get this parcel by id, testing for error and that new status of parcel no equal old status test parcel
	parcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, nextStatus, parcel.Status)
}

// TestGetByClient tests on getting all client parcels by client id
func TestGetByClient(t *testing.T) {
	// connection to db
	dbName := "tracker.db"
	db, err := sql.Open("sqlite", dbName)
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	// generate test parcels
	testParcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	testParcelsMap := map[int]Parcel{}
	// reset client number to same number for all test parcels
	client := randRange.Intn(9_999_999)
	for i := 0; i < len(testParcels); i++ {
		testParcels[i].Client = client
	}
	// add new parcels, testing for error and not empty id
	for i := 0; i < len(testParcels); i++ {
		id, err := store.Add(testParcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)
		// set test parcel id for each parcel
		testParcels[i].Number = id
		// saving each test parcel to parcel map with key as parcel id
		testParcelsMap[id] = testParcels[i]
	}
	// get all test parcels by client id saved on client variable,
	// testing for error and that number of getting test parcels equal number test parcels
	gettingTestParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(testParcels), len(gettingTestParcels))
	// testing for all parcels from gettingTestParcels variable is on testPracelsMap variable
	// and that getting test parcels equal test parcel
	for _, gettingTestParcel := range gettingTestParcels {
		testParcel, ok := testParcelsMap[gettingTestParcel.Number]
		require.True(t, ok)
		assert.Equal(t, testParcel, gettingTestParcel)
	}
}
