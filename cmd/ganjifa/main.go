package main

import (
	"math/rand"
	"os"
	"time"

	"ganjifa-web-sim/api"
	"ganjifa-web-sim/db"
	"ganjifa-web-sim/game"

	"github.com/jyotiskaghosh/ganjifa/cards"
	"github.com/jyotiskaghosh/ganjifa/match"

	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	rand.Seed(time.Now().UnixNano())

	logrus.Info("Starting..")

	for _, set := range cards.Sets {
		for uid, ctor := range *set {
			match.AddCard(uid, ctor)
		}
	}

	go game.GetLobby().StartTicker()

	api.CreateCardCache()

	db.Connect(os.Getenv("mongo_uri"), os.Getenv("mongo_name"))

	api.Start(os.Getenv("port"))

}
