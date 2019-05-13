package main

import (
	"flag"
	"github.com/fcortesgamez/go-mongodb-example/cmd/webshopd/config"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
)

func main() {

	flag.Parse()

	mongoCfg := config.MongoDBSettings{
		URL:    *config.MongoURL,
		User:   *config.MongoUser,
		Pass:   *config.MongoPass,
		AckMin: *config.MongoAckMin,
		FSync:  *config.MongoFSync,
	}

	log.WithFields(log.Fields{
		"URL":     mongoCfg.URL,
		"User":    mongoCfg.User,
		"Ack Min": mongoCfg.AckMin,
		"FSync":   mongoCfg.FSync,
	}).Info("Mongo DB settings")

	// MongoDB
	{
		mongo := config.MongoDB(&mongoCfg).DB("")
		defer mongo.Session.Close()
		mongo.Session.SetSafe(&mgo.Safe{W: mongoCfg.AckMin, FSync: mongoCfg.FSync})

		productSession := mongo.Session.Copy()
		productSession.SetMode(mgo.SecondaryPreferred, false)
	}

	log.Info("Done\n")
}
