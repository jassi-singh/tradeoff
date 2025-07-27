# TradeOff API Documentation

This document provides comprehensive documentation for the TradeOff game API, including REST endpoints, WebSocket communication, and data structures.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `[Your production URL]`

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Most endpoints require a valid access token in the Authorization header.

### Token Types

- **Access Token**: Short-lived token for API requests (typically 15 minutes)
- **Refresh Token**: Long-lived token for obtaining new access tokens (typically 7 days)

## REST API Endpoints

### Authentication

#### Login

Creates a new player session and returns authentication tokens.

```http
POST /api/login
Content-Type: application/json

{
  "username": "string"
}
```

**Response (200 OK):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "player": {
    "id": "uuid",
    "username": "string"
  }
}
```

#### Refresh Token

Refreshes an expired access token using a refresh token.

```http
POST /api/refresh
Authorization: Bearer <refresh_token>
```

**Response (200 OK):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Position Management

#### Create Position

Creates a new trading position for the authenticated player.

```http
POST /api/position
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "type": "long" | "short"
}
```

**Response (201 Created):**

```json
{
  "type": "long",
  "entryPrice": 45000.0,
  "entryTime": "2024-12-01T10:30:00Z",
  "quantity": 0.00222222,
  "pnl": 0.0,
  "pnlPercentage": 0.0
}
```

**Error Responses:**

- `400 Bad Request`: Invalid position type
- `401 Unauthorized`: Invalid or missing token
- `409 Conflict`: Player already has an active position
- `400 Bad Request`: Player has no balance

#### Close Position

Closes the player's active position and calculates final P&L.

```http
POST /api/close-position
Authorization: Bearer <access_token>
```

**Response (204 No Content):**

**Error Responses:**

- `401 Unauthorized`: Invalid or missing token
- `400 Bad Request`: No active position to close

## WebSocket API

### Connection

Establish a WebSocket connection for real-time game updates.

```http
GET /ws?token=<jwt_access_token>
```

The connection will be upgraded to WebSocket and the client will receive real-time game updates.

### Message Types

All WebSocket messages follow this structure:

```json
{
  "type": "message_type",
  "data": {
    /* message-specific data */
  }
}
```

#### Game State Sync

Sent when a player first connects, containing complete game state.

**Type:** `game_state_sync`

```json
{
  "type": "game_state_sync",
  "data": {
    "roundId": "uuid",
    "chartData": [...],
    "phase": "lobby" | "live" | "closed",
    "endTime": "2024-12-01T10:30:00Z",
    "balance": 100.0,
    "activePosition": null | {...},
    "closedPositions": [...],
    "pnl": 0.0,
    "activePnl": 0.0,
    "activePnlPercentage": 0.0,
    "longPositions": 5,
    "shortPositions": 3,
    "totalPlayers": 8
  }
}
```

#### New Round

Sent when a new round starts, resetting all player states.

**Type:** `new_round`

```json
{
  "type": "new_round",
  "data": {
    "roundId": "uuid",
    "chartData": [...],
    "phase": "lobby",
    "endTime": "2024-12-01T10:30:00Z",
    "balance": 100.0,
    "activePosition": null,
    "closedPositions": [],
    "pnl": 0.0,
    "activePnl": 0.0,
    "activePnlPercentage": 0.0,
    "longPositions": 0,
    "shortPositions": 0,
    "totalPlayers": 0
  }
}
```

#### Phase Update

Sent when the game phase changes.

**Type:** `phase_update`

```json
{
  "type": "phase_update",
  "data": {
    "phase": "lobby" | "live" | "closed",
    "endTime": "2024-12-01T10:30:00Z"
  }
}
```

#### Price Update

Sent with real-time price data updates for chart display.

**Type:** `price_update`

```json
{
  "type": "price_update",
  "data": {
    "priceData": {
      "time": 1701432000,
      "open": 45000.0,
      "high": 45100.0,
      "low": 44900.0,
      "close": 45050.0,
      "volume": 1000.0
    },
    "updateLast": true
  }
}
```

#### P&L Update

Sent to individual players with their current P&L information.

**Type:** `pnl_update`

```json
{
  "type": "pnl_update",
  "data": {
    "pnl": 25.0,
    "balance": 125.0,
    "activePnl": 25.0,
    "activePnlPercentage": 5.56
  }
}
```

#### Count Update

Sent with current player and position count statistics.

**Type:** `count_update`

```json
{
  "type": "count_update",
  "data": {
    "longPositions": 5,
    "shortPositions": 3,
    "totalPlayers": 8
  }
}
```

## Data Structures

### Player

```json
{
  "id": "uuid",
  "username": "string"
}
```

### Position

```json
{
  "type": "long" | "short",
  "entryPrice": 45000.0,
  "entryTime": "2024-12-01T10:30:00Z",
  "quantity": 0.00222222,
  "pnl": 25.0,
  "pnlPercentage": 5.56
}
```

### Closed Position

```json
{
  "type": "long" | "short",
  "entryPrice": 45000.0,
  "entryTime": "2024-12-01T10:30:00Z",
  "quantity": 0.00222222,
  "pnl": 25.0,
  "pnlPercentage": 5.56,
  "exitPrice": 45250.0,
  "exitTime": "2024-12-01T10:35:00Z"
}
```

### Price Data

```json
{
  "time": 1701432000,
  "open": 45000.0,
  "high": 45100.0,
  "low": 44900.0,
  "close": 45050.0,
  "volume": 1000.0
}
```

### Game Phase

- `lobby`: Waiting period before trading begins
- `live`: Active trading phase
- `closed`: Cooldown period after trading ends

### Position Type

- `long`: Betting that the price will go up
- `short`: Betting that the price will go down

## Error Handling

### HTTP Error Responses

All endpoints return appropriate HTTP status codes:

- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `204 No Content`: Request successful, no content to return
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or failed
- `409 Conflict`: Resource conflict (e.g., already has position)
- `500 Internal Server Error`: Server error

### Error Response Format

```json
{
  "error": "Error message",
  "status": 400
}
```

### WebSocket Error Handling

- **Connection Errors**: WebSocket connection will close with appropriate close codes
- **Authentication Errors**: Connection will be rejected if invalid token provided
- **Message Errors**: Invalid messages will be logged but won't crash the connection

## Rate Limiting

Currently, the API does not implement rate limiting. This will be added in future phases as the application scales.

## Security Considerations

1. **JWT Tokens**: Access tokens have short expiration times for security
2. **HTTPS**: Production deployments should use HTTPS
3. **Input Validation**: All inputs are validated on the server side
4. **SQL Injection**: Protected through parameterized queries
5. **XSS**: Frontend implements proper escaping and sanitization

## Game Rules

### Trading Mechanics

- **Starting Balance**: $100 USD per player per round
- **Position Limits**: One active position per player at a time
- **Full Investment**: Players invest their entire balance when creating a position
- **P&L Calculation**: Real-time based on current market price vs entry price

### Game Phases

- **Lobby (15s)**: Players join and wait for trading to begin
- **Live (60s)**: Active trading phase where positions can be created and closed
- **Cooldown (10s)**: Brief pause before next round begins

### Market Data

- **Asset**: Bitcoin/USD (X:BTCUSD)
- **Data Source**: Polygon.io API
- **Update Frequency**: Real-time during live phase
- **Historical Data**: Uses actual Bitcoin price history for authenticity

## Development Notes

### Testing the API

You can test the API using tools like curl, Postman, or any HTTP client:

```bash
# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser"}'

# Create position (using token from login)
curl -X POST http://localhost:8080/api/position \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"type": "long"}'
```

### WebSocket Testing

You can test WebSocket connections using tools like wscat:

```bash
# Install wscat
npm install -g wscat

# Connect to WebSocket
wscat -c "ws://localhost:8080/ws?token=<access_token>"
```

## Future Enhancements

The following features are planned for future phases:

1. **Leaderboards**: Player rankings and statistics
2. **Scheduled Rounds**: Pre-scheduled game sessions
3. **Social Features**: Player profiles and sharing
4. **Advanced Analytics**: Detailed trading statistics
5. **Mobile App**: Native mobile application
6. **Real-time Chat**: Player communication during games
