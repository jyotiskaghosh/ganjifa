package server

import "ganjifa-web-sim/db"

// Message is the default message struct
type Message struct {
	Header string `json:"header"`
}

// DecksMessage lists the users decks
type DecksMessage struct {
	Header string    `json:"header"`
	Decks  []db.Deck `json:"decks"`
}

// LobbyChatMessage is used to store chat messages
type LobbyChatMessage struct {
	Username  string `json:"username"`
	Color     string `json:"color"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

// LobbyChatMessages is used to store chat messages
type LobbyChatMessages struct {
	Header   string             `json:"header"`
	Messages []LobbyChatMessage `json:"messages"`
}

// UserMessage holds information about users
type UserMessage struct {
	Username    string   `json:"username"`
	Color       string   `json:"color"`
	Hub         string   `json:"hub"`
	Permissions []string `json:"permissions"`
}

// UserListMessage is used to send a list of online users
type UserListMessage struct {
	Header string        `json:"header"`
	Users  []UserMessage `json:"users"`
}

// MatchMessage holds information about a match
type MatchMessage struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Name string `json:"name"`
}

// MatchesListMessage is used to list open matches
type MatchesListMessage struct {
	Header  string         `json:"header"`
	Matches []MatchMessage `json:"matches"`
}
