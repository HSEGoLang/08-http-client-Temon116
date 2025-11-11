//go:build !solution

package cardgame

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://deckofcardsapi.com/api/deck"

type Client struct {
	baseURL string
	client  *http.Client
	output  io.Writer
}

func NewClient() *Client {
	return &Client{
		baseURL: defaultBaseURL,
		client:  http.DefaultClient,
		output:  nil,
	}
}

func (c *Client) PlayGame(userGuess int) (bool, error) {
	resp, err := c.client.Get(c.baseURL + "/new/shuffle/?deck_count=1")
	if err != nil {
		return false, fmt.Errorf("failed to create deck: %w", err)
	}
	defer resp.Body.Close()

	var deckResp DeckResponse
	if err := json.NewDecoder(resp.Body).Decode(&deckResp); err != nil {
		return false, fmt.Errorf("failed to decode deck response: %w", err)
	}

	deckID := deckResp.DeckID
	realCount := 0

	for {
		drawResp, err := c.client.Get(fmt.Sprintf("%s/%s/draw/?count=1", c.baseURL, deckID))
		if err != nil {
			return false, fmt.Errorf("failed to draw card: %w", err)
		}
		defer drawResp.Body.Close()

		var draw DrawResponse
		if err := json.NewDecoder(drawResp.Body).Decode(&draw); err != nil {
			return false, fmt.Errorf("failed to decode draw response: %w", err)
		}

		card := draw.Cards
		realCount++

		c.printf("%s of %s\n", card[0].Value, card[0].Suit)

		if card[0].Value == "QUEEN" {
			break
		}
	}
	if realCount == userGuess {
		c.printf("Вы угадали!\n")
		return true, nil
	} else {
		c.printf("Вы проиграли! Правильный ответ: %d\n", realCount)
		return false, nil
	}
}

func (c *Client) printf(format string, args ...interface{}) {
	if c.output != nil {
		fmt.Fprintf(c.output, format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}

func PlayGame(userGuess int) (bool, error) {
	return NewClient().PlayGame(userGuess)
}
