package match

import (
	"errors"

	"github.com/jyotiskaghosh/ganjifa/game-api/family"

	"github.com/sirupsen/logrus"
)

// GetCard returns a pointer to a Card by its ID from the given containers
func GetCard(id string, cards []*Card) (*Card, error) {
	for _, card := range cards {
		if card.id == id {
			return card, nil
		}
	}
	return nil, errors.New("Card was not found")
}

// AssertCardsIn returns true or false based on if the specified card ids are present in the source []*Card
func AssertCardsIn(src []*Card, test ...string) bool {
	for _, toTest := range test {
		ok := false

		for _, card := range src {
			if card.id == toTest {
				ok = true
			}
		}

		if !ok {
			return false
		}
	}

	return true
}

// Devolve ...
func Devolve(card *Card, src *Card) {
	cards := card.player.Search(
		Filter(card.Attachments(), func(x *Card) bool { return x.Family() != family.Equipment }),
		"Select a card",
		1,
		1,
		true)

	if len(cards) < 1 {
		card.player.match.Destroy(card, src)
		return
	}

	if err := cards[0].MoveCard(BATTLEZONE); err != nil {
		logrus.Debugf("Couldn't devolve: %s", err)
		return
	}

	cards[0].Tapped = card.Tapped

	for _, card := range card.Attachments() {
		card.AttachTo(cards[0])
	}

	// This is done to maintain a single identity for a creature
	card.id, cards[0].id = cards[0].id, card.id
	card.conditions, cards[0].conditions = cards[0].conditions, card.conditions

	card.player.match.Destroy(card, src)
}

// Filter filters a slice of cards according to the filter func
func Filter(cards []*Card, filter func(*Card) bool) []*Card {
	filtered := make([]*Card, 0)

	for _, mCard := range cards {
		if filter(mCard) {
			filtered = append(filtered, mCard)
		}
	}

	return filtered
}

// AllContainers returns an array of all containers
func AllContainers() []Container {
	return []Container{
		BATTLEZONE,
		SOUL,
		TRAPZONE,
		HAND,
		GRAVEYARD,
		DECK,
	}
}

// hideCards takes an array of *Card and returns an array of empty CardStates
func hideCards(cards *[]CardState) []CardState {
	arr := make([]CardState, 0)

	for _, c := range *cards {
		arr = append(arr, CardState{ID: c.ID})
	}

	return arr
}

// denormalizeCards takes an array of *Card and returns an array of CardState
func denormalizeCards(cards []*Card) []CardState {
	arr := make([]CardState, 0)

	for _, c := range cards {
		arr = append(arr, c.denormalizeCard())
	}

	return arr
}
