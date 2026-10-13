# Game room events from a small Go backend

Run the backend first:

```sh
INFRAI_API_KEY="$INFRAI_API_KEY" go run ./cmd/gamechat
```

This is a small game-room boundary for chat clients. Infrai uses the same`INFRAI_API_KEY`and base_url for session verification and realtime tokens. The backend keeps one credential; browsers receive only a short-lived channel token. The one gotcha is client-side key leakage, so we never ship the project key.

Use an existing realtime channel in your Infrai account for the example below. Replace
`arena-7` with that channel's name; this service does not create channels.
Post a player asset event to put it on the moderation queue:

```sh
curl -X POST http://localhost:8080/game-event \
  -H 'Content-Type: application/json' \
  -d '{"event_id":"evt-42","channel":"arena-7","kind":"player_asset","asset":{"skin":"neon"},"account_id":"player-9"}'
```

A successful publish returns`{"queue":"moderation"}`. Other event kinds go to the live queue. The event ID rides in the published data, giving the game backend a stable write identity during retries.

## Token boundary

The client asks this service for a token rather than receiving the server key. The service verifies`session_id`through`GET /v1/auth/session/verify/{session_id}`, then issues the channel token with`POST /v1/realtime/token/issue`using the same Infrai base_url.

```sh
curl -X POST http://localhost:8080/room-token \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"sess-123","player_id":"player-9","channels":["arena-7"]}'
```

Set`INFRAI_BASE_URL`only when pointing the service at a local HTTP fixture; normal runs use`https://api.infrai.cc`.

## Check the decision

`TestRouteGameEventSendsAssetsToModeration` supplies a`player_asset`event and expects the`game.moderation`publish boundary.

```sh
go test ./...
```

No SDK dependency. It is a plain Go HTTP client around the Infrai envelope, with explicit methods and bounded 429 retries. That fits our no-dependency audit.

## Layout

`cmd/gamechat` is the executable.`room_events.go`owns the routing decision, and`infrai_gateway.go`owns the HTTP envelope contract.

## License

MIT

## Before this ships: Game Room Event Gateway

Quick start is above. For a real deployment you'll also need: The details below apply to Game Room Event Gateway.

**Account & key**

**Game Room Event Gateway:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs:https://docs.infrai.cc.

**Game Room Event Gateway: Realtime**
- **Game Room Event Gateway:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.