package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"testing"
)

// testData is just a simple key/value pair intended to be persisted in the database unit tests.
type testData struct {
	Key   string
	Value string
}

func TestSetupAndReleaseDB(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f"}
	session, collections := SetupDB(t, namedCollections(names)...)

	for _, n := range names {
		_, exists := collectionMutexes[n]
		assert.True(t, exists, "Collection %s expected locked by mutex", n)

		// Insert some data
		data := []testData{
			{Key: "k1", Value: "v1"},
			{Key: "k2", Value: "v2"},
		}
		for _, d := range data {
			err := collections[n].Insert(d)
			assert.NoError(t, err, "No error expected inserting data %v in collection %s", d, n)
		}

		// Check the data by querying on the key
		for _, d := range data {
			var res testData
			err := collections[n].Find(bson.M{"key": d.Key}).One(&res)

			assert.NoError(t, err, "No error expected finding data using key %v in collection %s", d.Key, n)
			expectedData := testData{Key: res.Key, Value: res.Value}
			assert.Equal(t, expectedData, res, "Expected data %v found in collection %s", expectedData, n)
		}
	}

	ReleaseDB(t, session, collections)

	for _, n := range names {
		m, exists := collectionMutexes[n]

		assert.Nil(t, m, "Expected mutex for collection %s released", n)
		assert.False(t, exists, "Expected mutex for collection %s not present", n)
	}
}

// namedCollections returns a slice of NamedCollection for the given names.
func namedCollections(names []string) []NamedCollection {
	var collections []NamedCollection

	for _, n := range names {
		var collectionFunc func() *Collection
		collections = append(collections, NamedCollection{Name: n, CollectionFunc: &collectionFunc})
	}

	return collections
}

func TestResolveTestName(t *testing.T) {
	type testFunc func() string
	for i, fx := range []testFunc{a, b, c, d, e, f} {
		assert.Equal(t, "mongo_TestResolveTestName", fx(), fmt.Sprintf("%d. Expected resolved test name", i))
	}
}

func a() string {
	return ResolveTestName()
}

func b() string {
	return c()
}

func c() string {
	return d()
}

func d() string {
	return e()
}

func e() string {
	return f()
}

func f() string {
	return ResolveTestName()
}
