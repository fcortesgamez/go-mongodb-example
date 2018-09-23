package config

// MongoDBSettings represents all the Mongo DB settting needed by the application
type MongoDBSettings struct {
	URL  string
	User string
	Pass string
}
