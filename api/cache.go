package api

import (
	"sync"

	"github.com/jyotiskaghosh/ganjifa/game-api/cards"

	"github.com/sirupsen/logrus"
)

// CardInfo struct is used for the card database api
type CardInfo struct {
	UID          int    `json:"uid"`
	Name         string `json:"name"`
	Civilization string `json:"civilization"`
	Set          string `json:"set"`
}

// Register holds all the card info
var register []CardInfo = make([]CardInfo, 0)
var mutex *sync.Mutex = &sync.Mutex{}

// CreateCardCache loads all cards and creates a cache of the static data
func CreateCardCache() {

	for setID, set := range cards.Sets {

		for uid, c := range *set {

			card := c()

			register = append(register, CardInfo{
				UID:          uid,
				Name:         card.Name,
				Civilization: string(card.Civ),
				Set:          setID,
			})
		}
	}

	logrus.Infof("Loaded %v cards into the cache from %v sets", len(register), len(cards.Sets))

}

// GetCache returns a copy of the cache
func GetCache() []CardInfo {
	return register
}

// CacheHas returns true if the specified uid exist in the cache
func CacheHas(uid int) bool {

	mutex.Lock()

	defer mutex.Unlock()

	for _, c := range register {
		if c.UID == uid {
			return true
		}
	}

	return false

}
