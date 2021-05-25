package match

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/sirupsen/logrus"
)

// Container for cards
type Container string

// containers
const (
	DECK       Container = "deck"
	HAND       Container = "hand"
	GRAVEYARD  Container = "graveyard"
	BATTLEZONE Container = "battlezone"
	HIDDENZONE Container = "hiddenzone"
	SOUL       Container = "soul"
)

// LIFE is the default life of player
const LIFE int = 2000

// Writer is the interface for the output
type Writer interface {
	Write(interface{})
}

// PlayerReference ties a player to a Writer interface
type PlayerReference struct {
	Name   string
	Player *Player
	Writer Writer
}

// Write data
func (pr *PlayerReference) Write(msg interface{}) {
	pr.Writer.Write(msg)
}

// Player holds the player data
type Player struct {
	deck       []*Card
	hand       []*Card
	graveyard  []*Card
	battlezone []*Card
	hiddenzone []*Card
	soul       []*Card
	life       int

	mutex *sync.Mutex

	ready bool

	turn   bool
	turnNo int

	Action chan []string
	Cancel chan bool

	match *Match
}

// newPlayer returns a new player
func newPlayer(match *Match, turn bool) *Player {

	p := &Player{
		deck:       make([]*Card, 0),
		hand:       make([]*Card, 0),
		graveyard:  make([]*Card, 0),
		battlezone: make([]*Card, 0),
		hiddenzone: make([]*Card, 0),
		soul:       make([]*Card, 0),
		life:       LIFE,
		mutex:      &sync.Mutex{},
		Action:     make(chan []string),
		Cancel:     make(chan bool),
		turn:       turn,
		match:      match,
	}

	return p
}

// Name returns the username of the player
func (p *Player) Name() string {
	pr, err := p.match.playerRef(p)
	if err != nil {
		logrus.Debug(err)
		return "unknown"
	}

	return pr.Name
}

// IsPlayerTurn is it the Player's turnNo
func (p *Player) IsPlayerTurn() bool {
	return p.turn
}

// Turn returns turn no.
func (p *Player) Turn() int {
	return p.turnNo
}

// containerRef returns a pointer to one of the player's card zones based on the specified string
func (p *Player) containerRef(c Container) (*[]*Card, error) {

	switch c {
	case DECK:
		return &p.deck, nil
	case HAND:
		return &p.hand, nil
	case GRAVEYARD:
		return &p.graveyard, nil
	case BATTLEZONE:
		return &p.battlezone, nil
	case HIDDENZONE:
		return &p.hiddenzone, nil
	case SOUL:
		return &p.soul, nil
	default:
		return nil, errors.New("Invalid container")
	}
}

// Container returns a copy of one of the player's card zones based on the specified string
func (p *Player) Container(c Container) ([]*Card, error) {

	switch c {
	case DECK:
		return p.deck, nil
	case HAND:
		return p.hand, nil
	case GRAVEYARD:
		return p.graveyard, nil
	case BATTLEZONE:
		return p.battlezone, nil
	case HIDDENZONE:
		return p.hiddenzone, nil
	case SOUL:
		return p.soul, nil
	default:
		return nil, errors.New("Invalid container")
	}
}

// MapContainer performs the given action on all cards in the specified container
func (p *Player) MapContainer(containerName Container, fnc func(*Card)) {

	cards, err := p.Container(containerName)
	if err != nil {
		logrus.Debug(err)
		return
	}

	for _, card := range cards {
		fnc(card)
	}
}

// createDeck initializes a new deck from a list of card ids
func (p *Player) createDeck(cards []int) error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	deck := make([]*Card, 0)

	for _, card := range cards {

		c, err := NewCard(p, card)
		if err != nil {
			logrus.Warnf("Failed to create card with id %d", card)
			continue
		}

		deck = append(deck, c)
	}

	if len(deck) != 40 {
		return errors.New("deck must have exactly 40 cards")
	}

	count := make(map[int]int)

	for _, card := range deck {
		count[card.cardID]++
		if count[card.cardID] > 4 {
			return errors.New("deck must have only 4 copies of a card")
		}
	}

	p.deck = deck
	p.ready = true
	return nil
}

// ShuffleDeck randomizes the order of cards in the players deck
func (p *Player) ShuffleDeck() {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	rand.Shuffle(len(p.deck), func(i, j int) { p.deck[i], p.deck[j] = p.deck[j], p.deck[i] })
	p.match.Chat("Server", fmt.Sprintf("%s's deck was shuffled", p.Name()))
}

// PeekDeck returns references to the next n cards in the deck
func (p *Player) PeekDeck(n int) []*Card {

	result := make([]*Card, 0)

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.deck) < n {
		n = len(p.deck)
	}

	for i := 0; i < n; i++ {
		result = append(result, p.deck[i])
	}

	return result
}

// DrawCards moves n cards from the players deck to their hand
func (p *Player) DrawCards(n int) {

	if len(p.deck) < n {
		n = len(p.deck)
	}

	for i := 0; i < n; i++ {
		p.deck[i].MoveCard(HAND)
	}

	if n > 1 {
		p.match.Chat("Server", fmt.Sprintf("%s drew %v cards", p.Name(), n))
	} else {
		p.match.Chat("Server", fmt.Sprintf("%s drew %v card", p.Name(), n))
	}

	if len(p.deck) == 0 {
		p.match.End(p.match.Opponent(p), fmt.Sprintf("%s has no cards left in his deck", p.Name()))
	}
}

// HasCard checks if a container has a card
func (p *Player) HasCard(container Container, id string) bool {

	c, err := p.Container(container)
	if err != nil {
		logrus.Debug(err)
		return false
	}

	for _, card := range c {
		if card.ID == id {
			return true
		}
	}

	return false
}

// GetCard returns a pointer to a Card by its ID and container
func (p *Player) GetCard(id string) (*Card, error) {

	for _, card := range p.match.collectCards() {
		if card.ID == id {
			return card, nil
		}
	}

	return nil, errors.New("Card was not found")
}

// Damage reduces life of player
func (p *Player) Damage(source *Card, ctx *Context, health int) {

	if health <= 0 {
		return
	}

	ctx = NewContext(ctx.Match,
		&DamageEvent{
			Player: p,
			Source: source,
			Event:  ctx.Event,
			Health: health,
		})

	ctx.Override(func() {
		if event, ok := ctx.Event.(*DamageEvent); ok {
			p.mutex.Lock()
			p.life -= event.Health
			p.mutex.Unlock()
		}
	})

	ctx.ScheduleAfter(func() {
		ctx.Match.Chat("Server", fmt.Sprintf("%s did %d damage to %s", source.name, health, p.Name()))
	})

	ctx.Match.HandleFx(ctx)

	ctx.Match.BroadcastState()

	if p.life <= 0 {
		ctx.Match.End(p.match.Opponent(p), fmt.Sprintf("%s has no life left", p.Name()))
	}
}

// Heal reduces life of player
func (p *Player) Heal(source *Card, ctx *Context, health int) {

	if health <= 0 {
		return
	}

	ctx = NewContext(ctx.Match,
		&HealEvent{
			Player: p,
			Source: source,
			Health: health,
		})

	ctx.Override(func() {
		if event, ok := ctx.Event.(*HealEvent); ok {
			p.mutex.Lock()
			p.life += event.Health
			p.mutex.Unlock()
		}
	})

	ctx.ScheduleAfter(func() {
		ctx.Match.Chat("Server", fmt.Sprintf("%s healed %d life for %s", source.name, health, p.Name()))
	})

	ctx.Match.HandleFx(ctx)

	ctx.Match.BroadcastState()
}

// denormalized returns a server.PlayerState
func (p *Player) denormalized() PlayerState {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	state := PlayerState{
		Life:       p.life,
		Deck:       len(p.deck),
		Hand:       denormalizeCards(p.hand),
		Graveyard:  denormalizeCards(p.graveyard),
		Battlezone: denormalizeCards(p.battlezone),
		Hiddenzone: denormalizeCards(p.hiddenzone),
	}

	return state
}

// denormalizeCards takes an array of *Card and returns an array of CardState
func denormalizeCards(cards []*Card) []CardState {

	arr := make([]CardState, 0)

	for _, card := range cards {

		cs := CardState{
			ID:            card.ID,
			UID:           card.cardID,
			Name:          card.name,
			Civ:           card.civ,
			Tapped:        card.tapped,
			AttachedCards: denormalizeCards(card.Attachments()),
		}

		arr = append(arr, cs)
	}

	return arr
}

// hideCards takes an array of *Card and returns an array of empty CardStates
func hideCards(n int) []CardState {

	arr := make([]CardState, 0)

	for i := 0; i < n; i++ {

		cs := CardState{}

		arr = append(arr, cs)
	}

	return arr
}
