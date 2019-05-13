// Package mongo provides functionality to perform operations against a MongoDB database such us
// connecting/disconnecting to the DB (but never domain specific).
package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"sync"
)

const dbKey = "_db"

var (
	// counters keep track of the concurrent connections per database collection,
	// allowing for the detection of connection leaks.
	counters map[string]int64
	lock     sync.Mutex

	// MaxLeakConnection is the maximum allowed number of DB connections for a database handle or given collection name.
	MaxLeakConnection int64 = 2

	// Errors
	ErrConnectionLeak = errors.New("connection leak")
)

// DB represents a handle to the Mongo DB.
type DB struct {
	*mgo.Database
}

// Collection represents as the name suggests, a Mongo DB collection.
type Collection struct {
	*mgo.Collection
}

// CollectionIndex represents a Mongo DB collection index having the target collection being resolved at any desired
// time, as it is resolved by the given function Collection.
type CollectionIndex struct {
	Collection func() *Collection
	Index      mgo.Index
}

// CollectionIndices is just a slice of CollectionIndex.
type CollectionIndices []CollectionIndex

// init setup the initial state of the database connections per collection name.
func init() {
	counters = make(map[string]int64)
}

// DBSupplier returns a function which creates a named database from the given session s.
// It will create such by creating a copy of the given session s.
func DBSupplier(s *mgo.Session) func() (*DB, error) {
	return func() (*DB, error) {
		err := incr(dbKey)
		return &DB{s.Copy().DB("")}, err
	}
}

// Close closes the database session hold by the database db, releasing the hold tracked connection.
func (db *DB) Close() {
	decr(dbKey)
	db.Session.Close()
}

// CollectionSupplier returns a function which creates a named collection for the given collection c.
// It will create such by creating a copy of the given session s.
func CollectionSupplier(s *mgo.Session, c string) func() (*Collection, error) {
	return func() (*Collection, error) {
		err := incr(c)
		return &Collection{s.Copy().DB("").C(c)}, err
	}
}

// Close closes the database session hold by the collection c, releasing the hold tracked connection.
func (c *Collection) Close() {
	decr(c.Name)
	c.Database.Session.Close()
}

// EnsureIndex ensures an index with the given key exists, creating it when necessary.
func (ci *CollectionIndex) EnsureIndex() error {
	c := ci.Collection()
	defer c.Close()

	return c.EnsureIndex(ci.Index)
}

// incr increments a DB connection/collection tracked by the given collection name.
// It returns and error in case connections are detected to be leaking or nil otherwise.
func incr(name string) error {
	lock.Lock()

	var err error
	if counters[name]+1 < MaxLeakConnection {
		counters[name] += 1
	} else {
		err = ErrConnectionLeak
	}

	lock.Unlock()

	return err
}

// decr decrements a DB connection/collection tracked by the given collection name.
func decr(name string) {
	lock.Lock()
	if counters[name] > 0 {
		counters[name] -= 1
	}
	lock.Unlock()
}
