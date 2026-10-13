package gamechat

import (
	"context"
	"encoding/json"
	"errors"
)

type SessionVerifier interface {
	Call(context.Context, string, string, any, any) error
}

type TokenRequest struct {
	SessionID string   `json:"session_id"`
	PlayerID  string   `json:"player_id"`
	Channels  []string `json:"channels"`
}

type tokenIssueInput struct {
	ClientID     string   `json:"client_id"`
	Channels     []string `json:"channels"`
	Capabilities []string `json:"capabilities"`
	TTLSeconds   int      `json:"ttl_seconds"`
}

type tokenIssueOutput struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func IssueRoomToken(ctx context.Context, client SessionVerifier, request TokenRequest) (TokenResponse, error) {
	if request.SessionID == "" || request.PlayerID == "" || len(request.Channels) == 0 {
		return TokenResponse{}, errors.New("session_id, player_id, and channels are required")
	}
	var session json.RawMessage
	if err := client.Call(ctx, httpMethodGet, "/v1/auth/session/verify/"+request.SessionID, nil, &session); err != nil {
		return TokenResponse{}, err
	}
	var issued tokenIssueOutput
	input := tokenIssueInput{ClientID: request.PlayerID, Channels: request.Channels, Capabilities: []string{"subscribe", "publish"}, TTLSeconds: 900}
	if err := client.Call(ctx, httpMethodPost, "/v1/realtime/token/issue", input, &issued); err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{Token: issued.Token}, nil
}

const (
	httpMethodGet  = "GET"
	httpMethodPost = "POST"
)

type GameEvent struct {
	EventID   string          `json:"event_id"`
	Channel   string          `json:"channel"`
	Kind      string          `json:"kind"`
	Asset     json.RawMessage `json:"asset"`
	AccountID string          `json:"account_id"`
}

type publishInput struct {
	Channel   string `json:"channel"`
	Event     string `json:"event"`
	Data      string `json:"data"`
	AccountID string `json:"account_id"`
}

type PublishResult struct {
	Queue string `json:"queue"`
}

func RouteGameEvent(ctx context.Context, client SessionVerifier, event GameEvent) (PublishResult, error) {
	if event.EventID == "" || event.Channel == "" || event.Kind == "" || event.AccountID == "" {
		return PublishResult{}, errors.New("event_id, channel, kind, and account_id are required")
	}
	queue := "live"
	if event.Kind == "player_asset" {
		queue = "moderation"
	}
	payload, err := json.Marshal(struct {
		EventID string          `json:"event_id"`
		Kind    string          `json:"kind"`
		Queue   string          `json:"queue"`
		Asset   json.RawMessage `json:"asset"`
	}{EventID: event.EventID, Kind: event.Kind, Queue: queue, Asset: event.Asset})
	if err != nil {
		return PublishResult{}, err
	}
	input := publishInput{Channel: event.Channel, Event: "game." + queue, Data: string(payload), AccountID: event.AccountID}
	if err := client.Call(ctx, httpMethodPost, "/v1/realtime/publish", input, nil); err != nil {
		return PublishResult{}, err
	}
	return PublishResult{Queue: queue}, nil
}
