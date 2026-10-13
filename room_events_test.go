package gamechat

import (
	"context"
	"encoding/json"
	"testing"
)

type recordedCall struct {
	method string
	path   string
	input  any
}

type recordingClient struct {
	calls []recordedCall
}

func (c *recordingClient) Call(_ context.Context, method, path string, input any, output any) error {
	c.calls = append(c.calls, recordedCall{method: method, path: path, input: input})
	if value, ok := output.(*tokenIssueOutput); ok {
		value.Token = "channel-token"
	}
	return nil
}

func TestRouteGameEventSendsAssetsToModeration(t *testing.T) {
	client := &recordingClient{}
	result, err := RouteGameEvent(context.Background(), client, GameEvent{
		EventID: "evt-42", Channel: "arena-7", Kind: "player_asset", Asset: json.RawMessage(`{"skin":"neon"}`), AccountID: "player-9",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Queue != "moderation" {
		t.Fatalf("queue = %q, want moderation", result.Queue)
	}
	if len(client.calls) != 1 || client.calls[0].method != "POST" || client.calls[0].path != "/v1/realtime/publish" {
		t.Fatalf("unexpected calls: %#v", client.calls)
	}
	published := client.calls[0].input.(publishInput)
	if published.Event != "game.moderation" {
		t.Fatalf("event = %q, want game.moderation", published.Event)
	}
}
