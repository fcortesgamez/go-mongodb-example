package config

import "flag"

// Mongo DB settings
var (
	MongoURL  = flag.String("mongo.url", "localhost", "MongoDB connect string")
	MongoUser = flag.String("mongo.user", "", "MongoDB username")
	MongoPass = flag.String("mongo.pass*", "", "MongoDB password")
)
