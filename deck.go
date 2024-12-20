package masterduelmeta

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Deck struct {
	Author  any       `json:"author"`
	Created time.Time `json:"created"`
	Main    []struct {
		Card struct {
			ID   string `json:"_id"`
			Name string `json:"name"`
		} `json:"card"`
		Amount int `json:"amount"`
	} `json:"main"`
	Extra []struct {
		Card struct {
			ID   string `json:"_id"`
			Name string `json:"name"`
		} `json:"card"`
		Amount int `json:"amount"`
	} `json:"extra"`
	Side       []any  `json:"side"`
	URL        string `json:"url"`
	RankedType struct {
		Name           RankedType `json:"name"`
		ShortName      string     `json:"shortName"`
		Icon           string     `json:"icon"`
		StatsWeight    int        `json:"statsWeight"`
		Event          bool       `json:"event"`
		ActiveEvent    bool       `json:"activeEvent"`
		IncludeInStats bool       `json:"includeInStats"`
	} `json:"rankedType"`
	DeckType struct {
		Name string `json:"name"`
	} `json:"deckType"`
	TournamentType struct {
		Name        TournamentType `json:"name"`
		ShortName   string         `json:"shortName"`
		Icon        string         `json:"icon"`
		EnumSuffix  string         `json:"enumSuffix"`
		StatsWeight int            `json:"statsWeight"`
	} `json:"tournamentType"`
	TournamentPlacement  TournamentPlacement `json:"tournamentPlacement"`
	TournamentNumber     string              `json:"tournamentNumber"`
	CustomTournamentName any                 `json:"customTournamentName"`
	Engines              []struct {
		ID   string `json:"_id"`
		Name string `json:"name"`
	} `json:"engines"`
	SRPrice       int `json:"srPrice"`
	URPrice       int `json:"urPrice"`
	LinkedArticle struct {
		ID    string `json:"_id"`
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"linkedArticle"`
	Illegal bool `json:"illegal"`
}

// RankedType can be "Master I", "Win Streaks", "Master V", ..., or
// "Extra Zero Festival", "Top 10 Rating", ..., or
// empty string (probably decks from community tournaments)
type RankedType string

type Archetype string // e.g. "Blue-Eyes", "Kashtira"

type (
	TournamentType      string // e.g. "Duelist Cup"
	TournamentPlacement string // e.g. "Top 100"
)

func ParseDecks(data []byte) ([]Deck, error) {
	var decks []Deck
	err := json.Unmarshal(data, &decks)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}
	return decks, nil
}

func (d Deck) CheckIsNormalRank() bool {
	if CheckIsNormalRank(d.RankedType.Name) {
		return true
	}
	tournamentType := strings.ToLower(string(d.TournamentType.Name))
	// handle top 10 duelist cup decks does not have RankedType data
	return strings.Contains(tournamentType, "duelist cup")
}

// Archetype usually is the same as DeckType.Name,
// but sometimes this normalizes the deck name, e.g. "Swordsoul Tenyi" to "Swordsoul"
func (d Deck) Archetype() Archetype {
	archetype := NormalizeDeckTypeName(d.DeckType.Name)
	switch archetype {
	case Branded:
		// check if the deck has more than 2 "Red-Eyes Dark Dragoon"
		return Branded
	case SnakeEye, FireKing:
		// check if the deck uses "Legendary Fire King Ponix"
		if d.CheckContainsCard("Legendary Fire King Ponix") {
			return FireKing
		}
		return SnakeEye
	default:
		return archetype
	}
}

// CheckContainsCard checks if the deck contains a card with at least the specified amount.
// If minAmount is not provided, it defaults to 1.
// It searches through the Main, Extra, and Side cards of the deck.
func (d Deck) CheckContainsCard(cardName string, minAmount ...int) bool {
	amount := 1
	if len(minAmount) > 0 {
		amount = minAmount[0]
	}

	for _, card := range d.Main {
		if card.Card.Name == cardName && card.Amount >= amount {
			return true
		}
	}
	for _, card := range d.Extra {
		if card.Card.Name == cardName && card.Amount >= amount {
			return true
		}
	}
	return false
}
