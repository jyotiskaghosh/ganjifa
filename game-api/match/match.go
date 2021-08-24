package match

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Match struct
type Match struct {
	player1 *Player
	player2 *Player

	mutex *sync.Mutex

	started bool

	winner *Player
	quit   chan bool
}

// New returns a new match object
func New() *Match {
	return &Match{
		mutex: &sync.Mutex{},
		quit:  make(chan bool),
	}
}

// PlayerForWriter returns the player for a given writer or an error if the writer is not in  p1 or p2
func (m *Match) PlayerForWriter(w Writer) (*Player, error) {
	if m.player1.writer == w {
		return m.player1, nil
	}

	if m.player2.writer == w {
		return m.player2, nil
	}

	return nil, errors.New("not a player of this match")
}

// Started if the game has started
func (m *Match) Started() bool {
	return m.started
}

// Quit if the game has ended
func (m *Match) Quit() <-chan bool {
	return m.quit
}

// Winner returns true or false based on if the Player is the winner
func (m *Match) Winner(p *Player) bool {
	return p == m.winner
}

// AddPlayer adds a new player
func (m *Match) AddPlayer(name string, writer Writer) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch {
	case m.player1 == nil:
		m.player1 = newPlayer(name, writer, m, true)
	case m.player2 == nil:
		m.player2 = newPlayer(name, writer, m, false)
	default:
		return errors.New("players at max capacity")
	}

	return nil
}

// Parse processes the data provided by player
func (m *Match) Parse(w Writer, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from parsing message in match. %v", r)
		}
	}()

	p, err := m.PlayerForWriter(w)
	if err != nil {
		logrus.Debug(err)
		return
	}

	if p.wait {
		Warn(p, "Waiting for an action to resolve")
		return
	}

	if m.player1 != nil {
		m.player1.waiting(true)
		defer m.player1.waiting(false)
	}

	if m.player2 != nil {
		m.player2.waiting(true)
		defer m.player2.waiting(false)
	}

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	switch msg.Header {
	case "choose_deck":
		{
			if m.Started() {
				Warn(p, "match started cannot choose deck again")
				return
			}

			msg := CreateDeck{}
			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			if err := p.createDeck(msg.Cards); err != nil {
				Warn(p, err.Error())
				return
			}

			m.Start()
		}
	case "end_turn":
		{
			if m.started && p.turn {
				m.EndTurn()
			}
		}
	case "set_card":
		{
			if !m.started || !p.turn {
				return
			}

			var msg struct {
				ID string `json:"id"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			c, err := GetCard(msg.ID, p.CollectCards(HAND))
			if err != nil {
				logrus.Debug(err)
				return
			}

			if err := c.MoveCard(TRAPZONE); err != nil {
				logrus.Debug(err)
				return
			}
		}
	case "play_card":
		{
			if !m.started || !p.turn {
				return
			}

			var msg struct {
				ID string `json:"id"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			if ok := p.HasCard(msg.ID, HAND); ok {
				m.PlayCard(msg.ID)
			}
		}
	case "action":
		{
			var msg struct {
				Cards []string `json:"cards"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				Warn(p, "Invalid selection")
				return
			}

			m := make(map[string]bool)
			for _, card := range msg.Cards {
				m[card] = true
			}

			var cards []string
			for card := range m {
				cards = append(cards, card)
			}

			if p.Action != nil {
				p.Action <- cards
			}
		}
	case "cancel":
		{
			if p.Cancel != nil {
				p.Cancel <- true
			}
		}
	case "attack_player":
		{
			if !m.started || !p.turn {
				return
			}

			var msg struct {
				ID string `json:"id"`
			}

			if p.turnNo == 1 && p == m.player1 {
				Warn(p, "player 1 can't attack on first turn")
				return
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			m.AttackPlayer(p, msg.ID)
		}
	case "attack_creature":
		{
			if !m.started || !p.turn {
				return
			}

			var msg struct {
				ID string `json:"id"`
			}

			if p.turnNo == 1 && p == m.player1 {
				Warn(p, "player 1 can't attack on first turn")
				return
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			m.AttackCreature(p, msg.ID)
		}
	default:
		logrus.Debugf("Received message in incorrect format: %v", string(data))
	}
}

// Opponent returns the opponent of the given player
func (m *Match) Opponent(p *Player) *Player {
	if m.player1 == p {
		return m.player2
	}
	return m.player1
}

// CurrentPlayer returns the turn player
func (m *Match) CurrentPlayer() *Player {
	if m.player1.turn {
		return m.player1
	}
	return m.player2
}

// Chat sends a chat message
func (m *Match) Chat(sender string, message string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from panic during sending chat message. %v", r)
			return
		}
	}()

	msg := ChatMessage{
		Header:  "chat",
		Message: message,
		Sender:  sender,
	}

	m.player1.Write(msg)
	m.player2.Write(msg)
}

// MessagePlayer sends a chat message from the server to the specified player
func (m *Match) MessagePlayer(p *Player, message string) {
	m.WritePlayer(p, ChatMessage{
		Header:  "chat",
		Message: message,
		Sender:  "server",
	})
}

// WritePlayer sends a message to a player
func (m *Match) WritePlayer(p *Player, msg interface{}) {
	p.Write(msg)
}

// Warn sends a warning to the specified player ref
func Warn(p *Player, message string) {
	p.Write(ChatMessage{
		Header:  "warn",
		Message: message,
	})
}

// WarnPlayer sends a warning to the specified player
func (m *Match) WarnPlayer(p *Player, message string) {
	m.WritePlayer(p, ChatMessage{
		Header:  "warn",
		Message: message,
	})
}

// Battle handles a battle between two creatures
func (m *Match) Battle(attacker *Card, defender *Card, blocked bool) {
	if attacker.zone != BATTLEZONE || defender.zone != BATTLEZONE {
		return
	}

	ctx := NewContext(m, &Battle{Attacker: attacker, Defender: defender, Blocked: blocked})
	ctx.ScheduleAfter(func() {
		if attacker.GetAttack(ctx) > defender.GetDefence(ctx) {
			m.Destroy(defender, attacker)
		}
	})
	m.HandleFx(ctx)
}

// Destroy sends the given card to its players graveyard
func (m *Match) Destroy(card *Card, source *Card) {
	ctx := NewContext(m, &CreatureDestroyed{ID: card.id, Source: source})
	m.HandleFx(ctx)
}

// HandleFx ...
func (m *Match) HandleFx(ctx *Context) {
	defer m.BroadcastState()

	cards := make([]*Card, 0)

	// The player in which turn it is is to be handled first
	if m.player1.turn {
		cards = append(
			m.player1.CollectCards(AllContainers()...),
			m.player2.CollectCards(AllContainers()...)...,
		)
	} else {
		cards = append(
			m.player2.CollectCards(AllContainers()...),
			m.player1.CollectCards(AllContainers()...)...,
		)
	}

	for _, c := range cards {
		for _, h := range append(c.effects, c.conditions...) {
			if ctx.cancel {
				break
			}

			h(c, ctx)
		}
	}

	for _, h := range ctx.postFxs {
		if ctx.cancel {
			return
		}

		h()
	}
}

// BroadcastState sends the current game's state to both players, hiding the opponent's hand
func (m *Match) BroadcastState() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from panic during sending state update. %v", r)
			return
		}
	}()

	player1 := m.player1.denormalized()
	player2 := m.player2.denormalized()

	p1state := StateMessage{
		Header: "state_update",
		State: State{
			MyTurn:   m.player1.turn,
			Me:       player1,
			Opponent: player2,
		},
	}

	p2state := StateMessage{
		Header: "state_update",
		State: State{
			MyTurn:   m.player2.turn,
			Me:       player2,
			Opponent: player1,
		},
	}

	p1state.State.Opponent.Hand = hideCards(&p1state.State.Opponent.Hand)
	p1state.State.Opponent.Trapzone = hideCards(&p1state.State.Opponent.Trapzone)

	p2state.State.Opponent.Hand = hideCards(&p2state.State.Opponent.Hand)
	p2state.State.Opponent.Trapzone = hideCards(&p2state.State.Opponent.Trapzone)

	m.player1.Write(p1state)
	m.player2.Write(p2state)
}

// End ends the match
func (m *Match) End(winner *Player, reason string) {
	logrus.Debugf("Attempting to end match")

	m.Chat("server", fmt.Sprintf("%s won the match, %s", winner.Name(), reason))

	m.winner = winner

	m.quit <- true
	close(m.quit)
}

// NewAction prompts the user to make a selection of the specified []Cards
func (m *Match) NewAction(p *Player, cards []*Card, minSelections int, maxSelections int, text string, cancellable bool) {
	m.WritePlayer(p, ActionMessage{
		Header:        "action",
		Cards:         denormalizeCards(cards),
		Text:          text,
		MinSelections: minSelections,
		MaxSelections: maxSelections,
		Cancellable:   cancellable,
	})
}

// CloseAction closes the card selection popup for the given player
func (m *Match) CloseAction(p *Player) {
	m.WritePlayer(p, Message{
		Header: "close_action",
	})
}

// ShowCards shows the specified cards to the player with a message of why it is being shown
func (m *Match) ShowCards(p *Player, message string, cards []string) {
	m.WritePlayer(p, ShowCardsMessage{
		Header:  "show_cards",
		Message: message,
		Cards:   cards,
	})
}

// Highlight highlights creatures
func (m *Match) Highlight(ids ...string) {
	msg := HighlightMessage{"highlight", append(make([]string, 0), ids...)}
	m.player1.Write(msg)
	m.player2.Write(msg)
}

// changeCurrentPlayer changes the current player
func (m *Match) changeCurrentPlayer() {
	m.player1.turn = !m.player1.turn
	m.player2.turn = !m.player2.turn
}

// Start starts the match
func (m *Match) Start() {
	if m.started ||
		m.player1 == nil || m.player2 == nil ||
		!m.player1.ready || !m.player2.ready {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.started = true

	m.player1.ShuffleDeck()
	m.player2.ShuffleDeck()

	m.player1.DrawCards(5)
	m.player2.DrawCards(5)

	m.Chat("Server", "The match has begun!")

	// This is done to offset beginNewTurn which changes current player
	m.changeCurrentPlayer()

	m.beginNewTurn()
}

// beginNewTurn starts a new turn
func (m *Match) beginNewTurn() {
	m.changeCurrentPlayer()
	m.CurrentPlayer().turnNo++

	m.HandleFx(NewContext(m, &BeginTurnStep{}))

	m.untapStep()
}

// untapStep ...
func (m *Match) untapStep() {
	m.HandleFx(NewContext(m, &UntapStep{}))
	m.startOfTurnStep()
}

// startOfTurnStep ...
func (m *Match) startOfTurnStep() {
	m.HandleFx(NewContext(m, &StartOfTurnStep{}))
	m.drawStep()
}

// drawStep ...
func (m *Match) drawStep() {
	m.HandleFx(NewContext(m, &DrawStep{}))
	m.CurrentPlayer().DrawCards(1)
}

// endStep ...
func (m *Match) endStep() {
	m.HandleFx(NewContext(m, &EndStep{}))
	m.beginNewTurn()
}

// EndTurn is called when the player attempts to end their turn
// If the context is not cancelled by a card, the endStep is called
func (m *Match) EndTurn() {
	ctx := NewContext(m, &EndTurnEvent{})
	m.HandleFx(ctx)

	if !ctx.cancel {
		m.endStep()
	}
}

// PlayCard is called when the player attempts to play a card
func (m *Match) PlayCard(id string) {
	ctx := NewContext(m, &PlayCardEvent{
		ID: id,
	})
	m.HandleFx(ctx)
}

// Equip Fires a Equip event
func (m *Match) Equip(id string, target *Card) {
	m.HandleFx(NewContext(m, &Equip{
		ID:     id,
		Target: target,
	}))
}

// SpellCast Fires a SpellCast event
func (m *Match) SpellCast(id string, targets []*Card) {
	ids := []string{}

	for _, c := range targets {
		ids = append(ids, c.id)
	}

	m.Highlight(ids...)
	defer m.Highlight()

	m.HandleFx(NewContext(m, &SpellCast{
		ID:      id,
		Targets: targets,
	}))
}

// AttackPlayer is called when the player attempts to attack the opposing player
func (m *Match) AttackPlayer(p *Player, id string) {
	m.Highlight(id)
	defer m.Highlight()

	m.Opponent(p).waiting(false)

	ctx := NewContext(m, &AttackPlayer{
		ID: id,
	})
	m.HandleFx(ctx)

	// Tap card after attack
	// card gets tapped even if attack fails, this is done to prevent reusing before attacking effects
	card, err := GetCard(id, p.CollectCards(BATTLEZONE))
	if err != nil {
		logrus.Debug(err)
	} else {
		card.Tapped = true
	}

	m.BroadcastState()
}

// AttackCreature is called when the player attempts to attack an opponent's creature
func (m *Match) AttackCreature(p *Player, id string) {
	targets := p.Search(
		Filter(m.Opponent(p).CollectCards(BATTLEZONE), func(c *Card) bool { return c.Tapped }),
		"Select creature to attack",
		1,
		1,
		true,
	)

	if len(targets) > 0 {
		m.Highlight(id, targets[0].id)
		defer m.Highlight()
	}

	m.Opponent(p).waiting(false)

	ctx := NewContext(m, &AttackCreature{
		ID:       id,
		TargetID: targets[0].id,
	})
	m.HandleFx(ctx)

	// Tap card after attack
	// card gets tapped even if attack fails, this is done to prevent reusing before attacking effects
	card, err := GetCard(id, p.CollectCards(BATTLEZONE))
	if err != nil {
		logrus.Debug(err)
	} else {
		card.Tapped = true
	}

	m.BroadcastState()
}
