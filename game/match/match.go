package match

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/jyotiskaghosh/ganjifa/db"
	"github.com/jyotiskaghosh/ganjifa/server"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

var matches = make(map[string]*Match)
var matchesMutex = sync.Mutex{}

// Get returns a *Match from the specified id
func Get(id string) (*Match, error) {
	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	if m, ok := matches[id]; ok {
		return m, nil
	}

	return nil, errors.New("Not found")
}

var lobbyMatches = make(chan server.MatchesListMessage)

// Match struct
type Match struct {
	id        string
	matchName string
	host      string
	visible   bool

	match *match.Match

	created int64
	ending  bool
}

// Info struct
type Info struct {
	ID        string `json:"id"`
	MatchName string `json:"name"`
	Host      string `json:"host"`
	Visible   bool   `json:"visible"`
}

// Info returns match information in MatchInfo struct
func (m *Match) Info() Info {
	return Info{
		ID:        m.id,
		MatchName: m.matchName,
		Host:      m.host,
		Visible:   m.visible,
	}
}

// Matches returns a list of the current matches
func Matches() []string {
	result := make([]string, 0)

	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	for id := range matches {
		result = append(result, id)
	}

	return result
}

// New returns a new match object
func New(matchName string, host string, visible bool) *Match {
	id, err := shortid.Generate()

	if err != nil {
		id = uuid.New().String()
	}

	m := &Match{
		id:        id,
		matchName: matchName,
		host:      host,
		visible:   visible,
		match:     match.New(),

		created: time.Now().Unix(),
	}

	matchesMutex.Lock()

	matches[id] = m

	matchesMutex.Unlock()

	UpdateMatchList()

	go m.startTicker()

	logrus.Debugf("Created match %s", id)

	return m
}

// Name just returns "match", obligatory for a hub
func (m *Match) Name() string {
	return "match"
}

// LobbyMatchList returns the channel to receive match list updates
func LobbyMatchList() chan server.MatchesListMessage {
	return lobbyMatches
}

// UpdateMatchList sends a server.MatchesListMessage through the lobby channel
func UpdateMatchList() {
	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	matchesMessage := make([]server.MatchMessage, 0)

	for _, match := range matches {
		if !match.visible {
			continue
		}

		matchesMessage = append(matchesMessage, server.MatchMessage{
			ID:   match.id,
			Host: match.host,
			Name: match.matchName,
		})
	}

	update := server.MatchesListMessage{
		Header:  "matches",
		Matches: matchesMessage,
	}

	lobbyMatches <- update
}

func (m *Match) startTicker() {
	ticker := time.NewTicker(10 * time.Second) // tick every 10 seconds

	defer ticker.Stop()
	defer m.Dispose()
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from match ticker. %v", r)
		}
	}()

	for {
		select {
		case <-m.match.Quit():
			{
				logrus.Debugf("Closing match %s", m.id)
				m.ending = true
				return
			}
		case <-ticker.C:
			{
				// Close the match if it was not started within 10 minutes of creation
				if !m.match.Started() && m.created < time.Now().Unix()-60*10 {
					logrus.Debugf("Closing match %s", m.id)
					return
				}
			}
		}
	}
}

// Dispose closes the match, disconnects the clients and removes all references to it
func (m *Match) Dispose() {
	logrus.Debugf("Disposing match %s", m.id)

	defer func() {
		if r := recover(); r != nil {
			logrus.Warningf("Recovered from disposing a match. %v", r)
		}
	}()

	for _, s := range server.SocketsInHub(m.id) {
		s.Close()
	}

	matchesMutex.Lock()

	delete(matches, m.id)

	matchesMutex.Unlock()

	logrus.Debugf("Closed match with id %s", m.id)

	UpdateMatchList()
}

// Find returns a match with the specified id, or an error
func Find(id string) (*Match, error) {
	m := matches[id]

	if m != nil {
		return m, nil
	}

	return nil, errors.New("Match does not exist")
}

// Parse websocket messages
func (m *Match) Parse(s *server.Socket, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from parsing a message for a match. %v", r)
		}
	}()

	if m.ending {
		s.Write(match.ChatMessage{
			Header:  "warn",
			Message: "match has ended",
			Sender:  "server",
		})
		return
	}

	var message server.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return
	}

	switch message.Header {
	case "join_match":
		{
			if _, err := m.match.AddPlayer(s); err != nil {
				s.Write(match.ChatMessage{
					Header:  "warn",
					Message: err.Error(),
					Sender:  "server",
				})
				return
			}

			collection := db.Collection("decks")

			cur, err := collection.Find(context.TODO(), bson.M{
				"$or": []bson.M{
					{"owner": s.User.UID},
					{"standard": true},
				},
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			defer cur.Close(context.TODO())

			decks := make([]db.Deck, 0)

			for cur.Next(context.TODO()) {
				var deck db.Deck

				if err := cur.Decode(&deck); err != nil {
					continue
				}

				if deck.Owner == s.User.UID || deck.Standard {
					decks = append(decks, deck)
				}

				if deck.Owner == s.User.UID || deck.Standard {
					decks = append(decks, deck)
				}
			}

			s.Write(server.DecksMessage{
				Header: "choose_deck",
				Decks:  decks,
			})
		}
	case "chat":
		{
			var msg struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			m.match.Chat(s.User.Username, msg.Message)
		}
	default:
		{
			pr, err := m.match.PlayerForWriter(s)
			if err != nil {
				return
			}
			m.match.Parse(pr, data)
		}
	}
}

// OnSocketClose is called when a socket disconnects
func (m *Match) OnSocketClose(s *server.Socket) {
	if pr, err := m.match.PlayerForWriter(s); err == nil {
		match.Warn(pr, "Your opponent disconnected, the match will close soon.")
		m.match.End(m.match.Opponent(pr.Player), "opponent disconnected")
	}
}
