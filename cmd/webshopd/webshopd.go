package main

import (
	"flag"
	"fmt"

	"github.com/fcortesgamez/go-mongodb-example/cmd/webshopd/config"
)

func main() {

	flag.Parse()

	mongo := config.MongoDBSettings{
		URL:  *config.MongoURL,
		User: *config.MongoUser,
		Pass: *config.MongoPass,
	}

	fmt.Printf("Mongo DB settings: [URL: %s, User: %s]\n", mongo.URL, mongo.User)

	fmt.Printf("Done\n")
}
