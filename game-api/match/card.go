package match

import (
	"reflect"
	"sync"

	"github.com/jyotiskaghosh/ganjifa/game-api/civ"

	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

// Card stores card data
type Card struct {
	ID     string
	CardID int
	Player *Player
	Zone   Container

	Name    string
	Rank    int
	Civ     civ.Civilisation
	Family  string
	Attack  int
	Defence int

	Tapped bool

	handlers   []HandlerFunc
	AttachedTo *Card

	conditions []HandlerFunc // temporary handlers

	mutex *sync.Mutex
}

// NewCard returns a new, initialized card
func NewCard(p *Player, cardID int) (*Card, error) {

	c, err := CardCtor(cardID)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	id, err := shortid.Generate()
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	c.ID = id
	c.CardID = cardID
	c.Player = p
	c.Zone = DECK
	c.mutex = &sync.Mutex{}

	return c, nil
}

// Use allows different cards to hook into match events
// Can be compared to a typical middleware function
func (c *Card) Use(handlers ...HandlerFunc) {
	c.handlers = append(c.handlers, handlers...)
}

// Conditions returns all conditions
func (c *Card) Conditions() []HandlerFunc {
	return c.conditions
}

// AddCondition adds temporary handler functions
func (c *Card) AddCondition(handlers ...HandlerFunc) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conditions = append(c.conditions, handlers...)
}

// HasCondition returns true or false based on if an handler is added to the cards list of conditions
func (c *Card) HasCondition(handler HandlerFunc) bool {

	for _, condition := range c.conditions {
		if reflect.ValueOf(condition).Pointer() == reflect.ValueOf(handler).Pointer() {
			return true
		}
	}

	return false
}

// RemoveCondition removes all instances of the given handler from the cards conditions
func (c *Card) RemoveCondition(handler HandlerFunc) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	tmp := make([]HandlerFunc, 0)

	for _, condition := range c.conditions {

		if reflect.ValueOf(condition).Pointer() != reflect.ValueOf(handler).Pointer() {
			tmp = append(tmp, condition)
		}
	}

	c.conditions = tmp
}

// ClearConditions removes all conditions from the card
func (c *Card) ClearConditions() {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conditions = make([]HandlerFunc, 0)
}

// Attachments returns a copy of the card's attached cards
func (c *Card) Attachments() []*Card {

	result := make([]*Card, 0)

	cards, err := c.Player.Container(SOUL)
	if err != nil {
		logrus.Debug(err)
		return result
	}

	for _, card := range cards {
		if card.AttachedTo == c {
			result = append(result, card)
		}
	}

	return result
}

// Attach attaches card to c
func (c *Card) Attach(toAttach ...*Card) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, card := range toAttach {
		card.AttachedTo = c
		if err := card.MoveCard(SOUL); err != nil {
			logrus.Debug(err)
		}
	}
}

// MoveCard tries to move a card to container b
func (c *Card) MoveCard(destination Container) error {

	from, err := c.Player.ContainerRef(c.Zone)
	if err != nil {
		logrus.Debug(err)
		return err
	}

	to, err := c.Player.ContainerRef(destination)
	if err != nil {
		logrus.Debug(err)
		return err
	}

	c.Player.mutex.Lock()
	defer c.Player.mutex.Unlock()

	temp := make([]*Card, 0)

	for _, card := range *from {
		if card.ID != c.ID {
			temp = append(temp, card)
		}
	}

	*from = temp

	*to = append(*to, c)

	f := c.Zone
	c.Zone = destination

	c.Player.match.HandleFx(NewContext(c.Player.match, &CardMoved{
		ID:   c.ID,
		From: f,
		To:   destination,
	}))

	for _, card := range c.Attachments() {
		card.AttachedTo = c.AttachedTo
		card.Player.match.HandleFx(NewContext(card.Player.match, &CardMoved{
			ID:   card.ID,
			From: card.Zone,
			To:   destination,
		}))
	}

	return nil
}

// GetRank returns the rank of a given card
func (c *Card) GetRank(ctx *Context) int {

	e := &GetRankEvent{
		Card:  c,
		Event: ctx.Event,
		Rank:  c.Rank,
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Rank
}

// GetCivilisation returns the Civ of a given card
func (c *Card) GetCivilisation(ctx *Context) map[civ.Civilisation]bool {

	e := &GetCivilisationEvent{
		Card:  c,
		Event: ctx.Event,
		Civ:   map[civ.Civilisation]bool{c.Civ: true},
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Civ
}

// GetFamily returns the family of a given card
func (c *Card) GetFamily(ctx *Context) map[string]bool {

	e := &GetFamilyEvent{
		Card:   c,
		Event:  ctx.Event,
		Family: map[string]bool{c.Family: true},
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Family
}

// GetAttack returns the attack of a given card
func (c *Card) GetAttack(ctx *Context) int {

	e := &GetAttackEvent{
		Card:   c,
		Event:  ctx.Event,
		Attack: c.Attack,
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Attack
}

// GetDefence returns the defence of a given card
func (c *Card) GetDefence(ctx *Context) int {

	e := &GetDefenceEvent{
		Card:    c,
		Event:   ctx.Event,
		Defence: c.Defence,
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Defence
}
