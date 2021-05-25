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
	player1 *PlayerReference
	player2 *PlayerReference

	mutex *sync.Mutex

	started bool
	wait    bool // Wait prevents other moves from executing, used to prevent race conditions

	winner *Player
	Quit   chan bool
}

// New returns a new match object
func New() *Match {

	m := &Match{
		mutex: &sync.Mutex{},
		Quit:  make(chan bool),
	}

	return m
}

// playerRef returns the player ref for a given player
func (m *Match) playerRef(p *Player) (*PlayerReference, error) {

	if m.player1.Player == p {
		return m.player1, nil
	}

	if m.player2.Player == p {
		return m.player2, nil
	}

	return nil, errors.New("player is not a player of this match")
}

// PlayerForWriter returns the player ref for a given output or an error if the output is not p1 or p2
func (m *Match) PlayerForWriter(w Writer) (*PlayerReference, error) {

	if m.player1.Writer == w {
		return m.player1, nil
	}

	if m.player2.Writer == w {
		return m.player2, nil
	}

	return nil, errors.New("not a player of this match")
}

// Started if the game has started
func (m *Match) Started() bool {
	return m.started
}

// Winner returns true or false based on if the playerref is the winner
func (m *Match) Winner(pr *PlayerReference) bool {
	return pr.Player == m.winner
}

// AddPlayer adds a new player
func (m *Match) AddPlayer(w Writer) (*PlayerReference, error) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	var pr *PlayerReference

	switch {

	case m.player1 == nil:
		{
			pr = &PlayerReference{
				Name:   "player1",
				Player: newPlayer(m, true),
				Writer: w,
			}
			m.player1 = pr
		}

	case m.player2 == nil:
		{
			pr = &PlayerReference{
				Name:   "player2",
				Player: newPlayer(m, false),
				Writer: w,
			}
			m.player2 = pr
		}

	default:
		return nil, errors.New("players at max capacity")
	}

	return pr, nil
}

// Parse processes the data provided by player
func (m *Match) Parse(pr *PlayerReference, data []byte) {

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from parsing message in match. %v", r)
		}
	}()

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	switch msg.Header {

	case "choose_deck":
		{
			if m.started {
				Warn(pr, "match started cannot choose deck again")
				return
			}

			msg := CreateDeck{}
			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			if err := pr.Player.createDeck(msg.Cards); err != nil {
				Warn(pr, err.Error())
				return
			}

			m.Chat("Server", fmt.Sprintf("%s has chosen their deck", pr.Name))

			if m.player1 != nil && m.player2 != nil && m.player1.Player.ready && m.player2.Player.ready {
				m.start()
			}
		}

	case "end_turn":
		{
			if !pr.Player.turn {
				return
			}

			if m.wait {
				Warn(pr, "waiting for players to make an action")
				return
			}

			m.waiting(true)
			defer m.waiting(false)

			m.EndTurn()
		}

	case "set_card":
		{
			if !pr.Player.turn {
				return
			}

			if m.wait {
				Warn(pr, "waiting for players to make an action")
				return
			}

			m.waiting(true)
			defer m.waiting(false)

			var msg struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			c, err := pr.Player.GetCard(msg.ID)
			if err != nil {
				logrus.Debug(err)
				return
			}

			if err := c.MoveCard(HIDDENZONE); err != nil {
				logrus.Debug(err)
				return
			}

			m.BroadcastState()
		}

	case "play_card":
		{
			if !pr.Player.turn {
				return
			}

			if m.wait {
				Warn(pr, "waiting for players to make an action")
				return
			}

			m.waiting(true)
			defer m.waiting(false)

			var msg struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			if ok := pr.Player.HasCard(HAND, msg.ID); ok {
				m.PlayCard(msg.ID)
			}
		}

	case "action":
		{
			var msg struct {
				Cards []string `json:"cards"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				Warn(pr, "Invalid selection")
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

			pr.Player.Action <- cards
		}

	case "cancel":
		{
			pr.Player.Cancel <- true
		}

	case "attack_player":
		{
			if !pr.Player.turn {
				return
			}

			if m.wait {
				Warn(pr, "waiting for players to make an action")
				return
			}

			m.waiting(true)
			defer m.waiting(false)

			if pr.Player.turnNo == 1 && pr == m.player1 {
				Warn(pr, "player 1 can't attack on first turn")
				return
			}

			var msg struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			m.AttackPlayer(pr, msg.ID)
		}

	case "attack_creature":
		{
			if !pr.Player.turn {
				return
			}

			if m.wait {
				Warn(pr, "waiting for players to make an action")
				return
			}

			m.waiting(true)
			defer m.waiting(false)

			if pr.Player.turnNo == 1 && pr == m.player1 {
				Warn(pr, "player 1 can't attack on first turn")
				return
			}

			var msg struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			creatures := pr.Player.Search(
				m.Opponent(pr.Player).GetCreatures(),
				"Select 1 of your opponent's creature to attack",
				1,
				1,
				false)

			for _, c := range creatures {
				m.AttackCreature(pr, msg.ID, c.ID)
			}
		}

	default:
		{
			logrus.Debugf("Received message in incorrect format: %v", string(data))
		}
	}
}

// Opponent returns the opponent of the given player
func (m *Match) Opponent(p *Player) *Player {

	if m.player1.Player == p {
		return m.player2.Player
	}

	return m.player1.Player
}

// CurrentPlayer returns the turn player
func (m *Match) CurrentPlayer() *PlayerReference {

	if m.player1.Player.turn {
		return m.player1
	}

	return m.player2
}

// Chat sends a chat message with color
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

// WritePlayer sends a message to a player
func (m *Match) WritePlayer(p *Player, msg interface{}) {

	pr, err := m.playerRef(p)
	if err != nil {
		logrus.Debug(err)
		return
	}

	pr.Write(msg)
}

// Warn sends a warning to the specified player ref
func Warn(p *PlayerReference, message string) {

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

	ctx.Override(func() {
		if attacker.GetAttack(ctx) > defender.GetDefence(ctx) {
			m.Destroy(defender, attacker, fmt.Sprintf("%s was destroyed by %s", defender.name, attacker.name))
		}
	})

	m.HandleFx(ctx)

	m.BroadcastState()
}

// Destroy sends the given card to its players graveyard
func (m *Match) Destroy(card *Card, source *Card, text string) {

	ctx := NewContext(m, &CreatureDestroyed{Card: card, Source: source})

	ctx.ScheduleAfter(func() {
		m.Chat("Server", text)
	})

	m.HandleFx(ctx)
}

func (m *Match) collectCards() []*Card {

	players := make([]*PlayerReference, 0)

	// The player in which turn it is is to be handled first
	if m.player1.Player.turn {
		players = append(players, m.player1, m.player2)
	} else {
		players = append(players, m.player2, m.player1)
	}

	cards := make([]*Card, 0)

	for _, p := range players {

		cards = append(cards, p.Player.battlezone...)
		cards = append(cards, p.Player.soul...)
		cards = append(cards, p.Player.hiddenzone...)
		cards = append(cards, p.Player.hand...)
		cards = append(cards, p.Player.graveyard...)
		cards = append(cards, p.Player.deck...)
	}

	return cards
}

// ResolveEvent runs a check on a particular context
func (m *Match) ResolveEvent(ctx *Context) {

	for _, card := range m.collectCards() {

		for _, h := range append(card.handlers, card.conditions...) {

			if ctx.cancel {
				return
			}

			h(card, ctx)
		}
	}
}

// HandleFx ...
func (m *Match) HandleFx(ctx *Context) {

	m.ResolveEvent(ctx)

	for _, h := range ctx.preFxs {

		if ctx.cancel {
			return
		}

		h()
	}

	if ctx.cancel {
		return
	}

	if ctx.mainFx != nil {

		ctx.mainFx()
	}

	if ctx.cancel {
		return
	}

	for _, h := range ctx.postFxs {

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

	player1 := m.player1.Player.denormalized()
	player2 := m.player2.Player.denormalized()

	p1state := StateMessage{
		Header: "state_update",
		State: State{
			MyTurn:   m.player1.Player.turn,
			Me:       player1,
			Opponent: player2,
		},
	}

	p2state := StateMessage{
		Header: "state_update",
		State: State{
			MyTurn:   m.player2.Player.turn,
			Me:       player2,
			Opponent: player1,
		},
	}

	p1state.State.Opponent.Hand = hideCards(len(p1state.State.Opponent.Hand))
	p1state.State.Opponent.Hiddenzone = hideCards(len(p1state.State.Opponent.Hiddenzone))

	p2state.State.Opponent.Hand = hideCards(len(p2state.State.Opponent.Hand))
	p2state.State.Opponent.Hiddenzone = hideCards(len(p2state.State.Opponent.Hiddenzone))

	m.player1.Write(p1state)
	m.player2.Write(p2state)
}

// End ends the match
func (m *Match) End(winner *Player, reason string) {

	logrus.Debugf("Attempting to end match")

	m.Chat("server", fmt.Sprintf("%s won the match, %s", winner.Name(), reason))

	m.winner = winner
	m.Quit <- true
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

// changeCurrentPlayer changes the current player
func (m *Match) changeCurrentPlayer() {

	m.player1.Player.mutex.Lock()
	m.player1.Player.turn = !m.player1.Player.turn
	m.player1.Player.mutex.Unlock()

	m.player2.Player.mutex.Lock()
	m.player2.Player.turn = !m.player2.Player.turn
	m.player2.Player.mutex.Unlock()
}

// start starts the match
func (m *Match) start() {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.started = true

	m.player1.Player.ShuffleDeck()
	m.player2.Player.ShuffleDeck()

	m.player1.Player.DrawCards(5)
	m.player2.Player.DrawCards(5)

	m.Chat("Server", "The match has begun!")

	// This is done to offset beginNewTurn which changes current player
	m.changeCurrentPlayer()

	m.beginNewTurn()
}

// waiting assigns bool value to m.wait
func (m *Match) waiting(wait bool) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.wait = wait
}

// beginNewTurn starts a new turn
func (m *Match) beginNewTurn() {

	m.changeCurrentPlayer()
	m.CurrentPlayer().Player.turnNo++

	m.HandleFx(NewContext(m, &BeginTurnStep{}))

	m.BroadcastState()

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

	m.Chat("Server", fmt.Sprintf("Your turn, %s", m.CurrentPlayer().Player.Name()))

	m.drawStep()
}

// drawStep ...
func (m *Match) drawStep() {

	m.HandleFx(NewContext(m, &DrawStep{}))

	p := m.CurrentPlayer().Player
	p.DrawCards(1)

	m.BroadcastState()
}

// endStep ...
func (m *Match) endStep() {

	m.HandleFx(NewContext(m, &EndStep{}))

	m.Chat("Server", fmt.Sprintf("%s ended their turn", m.CurrentPlayer().Player.Name()))

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

	m.HandleFx(NewContext(m, &PlayCardEvent{
		ID: id,
	}))

	m.BroadcastState()
}

// AttackPlayer is called when the player attempts to attack the opposing player
func (m *Match) AttackPlayer(pr *PlayerReference, id string) {

	ctx := NewContext(m, &AttackPlayer{
		ID: id,
	})

	m.HandleFx(ctx)

	m.BroadcastState()
}

// AttackCreature is called when the player attempts to attack the opposing player
func (m *Match) AttackCreature(pr *PlayerReference, id string, targetID string) {

	ctx := NewContext(m, &AttackCreature{
		ID:       id,
		TargetID: targetID,
	})

	m.HandleFx(ctx)

	m.BroadcastState()
}

// GetCard returns the *Card from a container from either players
func (m *Match) GetCard(id string) (*Card, error) {

	c, err := m.player1.Player.GetCard(id)
	if err != nil {
		return m.player2.Player.GetCard(id)
	}

	return c, nil
}
