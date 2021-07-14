package match

import (
	"errors"
	"strconv"
)

// CardConstructor initializes and returns a card
type CardConstructor func() *Card

var ctors = make(map[int]CardConstructor)

// AddCard adds a new card constructor to ctors
func AddCard(id int, ctor CardConstructor) {
	ctors[id] = ctor
}

// CardCtor returns a *Card from cardID, or an error if it does not exist
func CardCtor(id int) (*Card, error) {
	if ctors[id] == nil {
		return nil, errors.New("Card ctor does not exist for id " + strconv.Itoa(id))
	}

	return ctors[id](), nil
}
