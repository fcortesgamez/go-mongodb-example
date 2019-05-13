package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// NamedCollection represents a database collection represented as a pair (Name/CollectionFunc), where Name is just
// the collection name, and CollectionFunc is a function to resolve a handle to such database collection.
type NamedCollection struct {
	Name           string
	CollectionFunc *func() *Collection
}

var mainMutex = &sync.Mutex{}
var collectionMutexes = map[string]*sync.Mutex{}

// SetupDB is a function which allows allocating a given set of database collections in unit tests.
// Such collections will be allocated on a database which name will be guarantied to be unique, as it is inferred from
// the test's package and file name.
// It will create a connection to Mongo DB on localhost.
func SetupDB(t *testing.T, namedCollections ...NamedCollection) (*mgo.Session, map[string]*Collection) {
	nameDB := ResolveTestName()
	collections := map[string]*Collection{}

	uri := fmt.Sprintf("mongodb://localhost/%s", nameDB)
	var err error
	session, err := mgo.Dial(uri)
	if err != nil {
		t.Fatalf("Failed to connect to Mongo on %q. Error: %v", uri, err)
	}
	session.SetSafe(&mgo.Safe{W: 1})

	// Assign all the DB collection handles
	for _, c := range namedCollections {
		mainMutex.Lock()
		_, exists := collectionMutexes[c.Name]
		if !exists {
			collectionMutexes[c.Name] = &sync.Mutex{}
		}
		collectionMutexes[c.Name].Lock()
		mainMutex.Unlock()

		*c.CollectionFunc = collectionSupplier(t, session.DB(nameDB).Session, nameDB, c.Name)
		f := *c.CollectionFunc
		collections[c.Name] = f()
	}

	return session, collections
}

// ReleaseDB will make sure the allocated database session and collection handles are released.
func ReleaseDB(t *testing.T, session *mgo.Session, collections map[string]*Collection) {
	released, failed := releaseDB(session, collections)

	for _, msg := range released {
		t.Logf(msg)
	}
	for _, msg := range failed {
		t.Errorf(msg)
	}
}

func releaseDB(session *mgo.Session, collections map[string]*Collection) ([]string, []string) {
	defer session.Close()
	defer session.LogoutAll()

	// Release the database handles
	namesDB := make(map[string]struct{})
	var releaseDBCollection = func(c *Collection) {
		defer collectionMutexes[c.Name].Unlock()
		defer c.Close()
		namesDB[c.Database.Name] = struct{}{}
		delete(collectionMutexes, c.Name)
	}

	for _, c := range collections {
		releaseDBCollection(c)
	}

	// Release the allocated test database (full database drop)
	var released []string
	var failed []string
	for name := range namesDB {
		err := session.DB(name).DropDatabase()
		if err != nil {
			failed = append(failed, fmt.Sprintf("DB %s. Failed to release. Error: %s", name, err))
		} else {
			released = append(released, fmt.Sprintf("DB %s. Successfully released", name))
		}
	}
	return released, failed
}

// ResolveTestName resolves the test name no matter at which position in the stack this method is after the test method
// has called it via any number of stacked called functions.
func ResolveTestName() string {
	pc := make([]uintptr, 10)
	runtime.Callers(1, pc)

	var filteredPc []uintptr
	for i := len(pc) - 1; i >= 0; i-- {
		if pc[i] != 0 {
			filteredPc = append(filteredPc, pc[i])
		}
	}

	f := runtime.FuncForPC(filteredPc[2])

	// Note: Only take the package and test function. For instance: with a function 'a' in a package 'p'
	// that would resolve to "p_a" (replacing the dot by underscore)
	lastSlashIdx := strings.LastIndex(f.Name(), "/")
	fullName := f.Name()
	testName := fullName[lastSlashIdx+1:]
	return strings.Replace(testName, ".", "_", -1)
}

func collectionSupplier(t *testing.T, session *mgo.Session, dbName, collection string) func() *Collection {
	return func() *Collection {
		if err := incr(collection); err != nil {
			t.Fatalf(err.Error())
		}

		return &Collection{session.Copy().DB(dbName).C(collection)}
	}
}
