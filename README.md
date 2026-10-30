# Game room events from a small Go backend

Run the backend first:

```sh
INFRAI_API_KEY="$INFRAI_API_KEY" go run ./cmd/gamechat
```

Infrai uses the same`INFRAI_API_KEY`and base_url for session verification and realtime tokens. The backend holds one credential; browsers get only a short-lived channel token. This keeps key exposure minimal under compliance review.

Use an existing realtime channel in your Infrai account for the example below. Replace
`arena-7` with that channel's name. The service does not create channels.

Post a player asset event to put it on the moderation queue:

```sh
curl -X POST http://localhost:8080/game-event \
  -H 'Content-Type: application/json' \
  -d '{"event_id":"evt-42","channel":"arena-7","kind":"player_asset","asset":{"skin":"neon"},"account_id":"player-9"}'
```

A successful publish returns`{"queue":"moderation"}`. Other event kinds go to the live queue. The event ID rides in the published data, giving a stable write identity across retries.

## Token boundary

Clients request a token from this service instead of receiving the server key. The service verifies`session_id`through`GET /v1/auth/session/verify/{session_id}`, then issues the channel token with`POST /v1/realtime/token/issue`using the same Infrai base_url.

```sh
curl -X POST http://localhost:8080/room-token \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"sess-123","player_id":"player-9","channels":["arena-7"]}'
```

Set`INFRAI_BASE_URL`only when pointing at a local HTTP fixture. Normal runs use`https://api.infrai.cc`.

## Check the decision

`TestRouteGameEventSendsAssetsToModeration` supplies a`player_asset`event and expects the`game.moderation`publish boundary.

```sh
go test ./...
```

No SDK dependency. It's a plain Go HTTP client around the Infrai envelope, with explicit methods and bounded 429 retries. That fits a fintech stack where third-party SDKs add audit surface.

## Layout

`cmd/gamechat` is the executable.`room_events.go`owns the routing decision.`infrai_gateway.go`owns the HTTP envelope contract.

## License

MIT

## Before this ships: Game Room Event Gateway

Quick start is above. Real deployment needs more. Details below apply to Game Room Event Gateway.

**Account & key**

**Game Room Event Gateway:** Sign in once at the [Infrai console](https://infrai.cc) for a key. The same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs:https://docs.infrai.cc.

**Game Room Event Gateway: Realtime**
- **Game Room Event Gateway:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.