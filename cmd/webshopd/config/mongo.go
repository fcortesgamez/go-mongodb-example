package config

import (
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
)

// MongoDBSettings represents all the Mongo DB setting needed by the application
type MongoDBSettings struct {
	URL    string
	User   string
	Pass   string
	AckMin int
	FSync  bool
}

// MongoDB will return a connection with a Mongo DB by the given settings cfg.
func MongoDB(cfg *MongoDBSettings) *mgo.Session {
	// TODO: It would be better to return the session and and error, so that logging can be done from the application side
	s, err := mgo.Dial(cfg.URL)
	if err != nil {
		log.WithFields(log.Fields{
			"URL":   cfg.URL,
			"error": err,
		}).Error("Failed to connect to Mongo DB")
	}

	if cfg.User != "" {
		creds := &mgo.Credential{Username: cfg.User, Password: cfg.Pass}
		if err := s.Login(creds); err != nil {
			s.Close()
			log.WithFields(log.Fields{
				"URL":   cfg.URL,
				"user":  cfg.User,
				"error": err,
			})
		}
	}
	return s
}
