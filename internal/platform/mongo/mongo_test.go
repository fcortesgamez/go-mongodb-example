package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBSupplierLeakConnections(t *testing.T) {
	s := dbSession()
	defer s.Close()

	dbSupplier := DBSupplier(s)

	// Supply a DB handle via a first session copy (no leak expected)
	db, err := dbSupplier()
	assert.NoError(t, err, "No errors expected")
	assert.NotNil(t, db, "Database handle copy expected")
	assert.Equal(t, int64(1), counters[dbKey], "Expected concurrent connections")

	// Supply a DB handle via a second session copy (leak expected as there are 2 session copies)
	_, err = dbSupplier()
	assert.Error(t, err, "Error expected %+v", ErrConnectionLeak)
	assert.Equal(t, int64(1), counters[dbKey], "Expected concurrent connections")

	// Close the leaked DB connection, to make sure hold connections are released
	db.Close()
	assert.Equal(t, int64(0), counters[dbKey], "Expected concurrent connections")

	// Close the DB connection again, which should not decreased the hold database connections (already 0)
	db.Close()
	assert.Equal(t, int64(0), counters[dbKey], "Expected concurrent connections")
}

func TestCollectionSupplierLeakConnections(t *testing.T) {
	s := dbSession()
	defer s.Close()

	collectionSupplier := CollectionSupplier(s, "c1")

	// Supply a collection handle via a first session copy (no leak expected)
	c, err := collectionSupplier()
	assert.NoError(t, err, "No errors expected")
	assert.NotNil(t, c, "Collection handle copy expected")
	assert.Equal(t, int64(1), counters["c1"], "Expected concurrent connections")

	// Supply a collection handle via a second session copy (leak expected as there are 2 session copies)
	_, err = collectionSupplier()
	assert.Error(t, err, "Error expected %+v", ErrConnectionLeak)
	assert.Equal(t, int64(1), counters["c1"], "Expected concurrent connections")

	// Close the leaked DB connection, to make sure hold connections are released
	c.Close()
	assert.Equal(t, int64(0), counters["c1"], "Expected concurrent connections")

	// Close the DB connection again, which should not decreased the hold collection connections (already 0)
	c.Close()
	assert.Equal(t, int64(0), counters["c1"], "Expected concurrent connections")
}

// dbSession connects to a Mongo DB on localhost, returning the Mongo connected session.
func dbSession() *mgo.Session {
	uri := "mongodb://localhost/"

	s, err := mgo.Dial(uri)
	if err != nil {
		panic(fmt.Sprintf("Connect to %q: %s", uri, err))
	}

	s.SetSafe(&mgo.Safe{W: 1})

	return s
}
