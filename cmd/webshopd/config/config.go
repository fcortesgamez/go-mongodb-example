package config

import "flag"

// Mongo DB settings
var (
	MongoURL  = flag.String("mongo.url", "localhost", "MongoDB connect string")
	MongoUser = flag.String("mongo.user", "", "MongoDB username")
	// TODO: Maybe I can keep it simple for now.
	MongoPass   = flag.String("mongo.pass*", "", "MongoDB password")
	MongoAckMin = flag.Int("mongo.write.ack-min", 1, "Number of MongoDB servers to acknoledge for a confirmation")
	MongoFSync  = flag.Bool("mongo.write.fsync", true, "Synchronize filesystem on MongoDB servers for a confirmation")
)
