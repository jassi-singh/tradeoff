# TradeOff Architecture Overview

This document provides a comprehensive overview of the TradeOff system architecture, explaining the current implementation, data flow, and technical decisions made during Phase 1 development.

## System Overview

TradeOff is a real-time, multiplayer stock market simulator built as a monolithic application with clear separation of concerns. The system consists of:

- **Backend**: Go-based API server with WebSocket support
- **Frontend**: Next.js React application with real-time UI
- **Database**: PostgreSQL for persistent data storage
- **Infrastructure**: Docker containerization for consistent deployment

## High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   Database      │
│   (Next.js)     │◄──►│    (Go)         │◄──►│  (PostgreSQL)   │
│                 │    │                 │    │                 │
│ • React UI      │    │ • REST API      │    │ • Player Data   │
│ • WebSocket     │    │ • WebSocket Hub │    │ • Auth Tokens   │
│ • State Mgmt    │    │ • Game Logic    │    │                 │
│ • Real-time     │    │ • Market Data   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │
         │                       │
         └───────────────────────┘
              WebSocket
              Real-time Updates
```

## Backend Architecture

### Service Layer Design

The backend follows a clean architecture pattern with distinct layers:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP/WebSocket Layer                    │
├─────────────────────────────────────────────────────────────┤
│  Handlers (auth_handler.go, position_handler.go, etc.)    │
├─────────────────────────────────────────────────────────────┤
│                    Service Layer                           │
├─────────────────────────────────────────────────────────────┤
│  RoundManager  │  PlayerService  │  MarketService  │  Hub │
├─────────────────────────────────────────────────────────────┤
│                    Storage Layer                           │
├─────────────────────────────────────────────────────────────┤
│                    PostgreSQL                              │
└─────────────────────────────────────────────────────────────┘
```

### Core Services

#### 1. RoundManager

**Purpose**: Orchestrates the entire game lifecycle

**Key Responsibilities**:

- Manages game phases (Lobby → Live → Cooldown)
- Coordinates market data loading and streaming
- Broadcasts game state updates to all players
- Handles round transitions and player state resets

**Key Methods**:

```go
func (r *RoundManager) Run()                    // Main game loop
func (r *RoundManager) transitionToLive()       // Phase transitions
func (r *RoundManager) sendPriceUpdate()        // Real-time price updates
func (r *RoundManager) GetGameState()           // Player state sync
```

#### 2. PlayerService

**Purpose**: Manages all player sessions and trading logic

**Key Responsibilities**:

- Thread-safe player session management
- Position creation and closing logic
- Real-time P&L calculations
- Player statistics aggregation

**Key Methods**:

```go
func (s *PlayerService) CreatePosition()        // Position creation
func (s *PlayerService) ClosePosition()         // Position closing
func (s *PlayerService) GetPlayerStat()         // P&L calculations
func (s *PlayerService) UpdateAllPlayerPnl()    // Real-time updates
```

#### 3. Hub

**Purpose**: Manages WebSocket connections and message routing

**Key Responsibilities**:

- WebSocket connection lifecycle management
- Message broadcasting to all connected clients
- Direct message routing to specific players
- Connection cleanup and error handling

**Key Methods**:

```go
func (h *Hub) Run()                            // Message routing loop
func (h *Hub) Broadcast()                      // Broadcast to all clients
func (h *Hub) SendDirect()                     // Direct message to player
```

#### 4. MarketService

**Purpose**: Fetches and manages market data

**Key Responsibilities**:

- Integration with Polygon.io API
- Historical data loading for game rounds
- Data transformation for chart display

#### 5. Leaderboard System

**Purpose**: Manages real-time player rankings and statistics

**Key Responsibilities**:

- Calculates active balances including unrealized P&L
- Sorts players by performance for leaderboard display
- Provides real-time updates to all connected clients
- Limits leaderboard to top 20 players for performance

**Key Methods**:

```go
func (s *PlayerService) GetLeaderboard() []domain.LeaderboardPlayer
```

### Data Flow

#### Game Initialization

1. **Server Start**: RoundManager initializes with empty game state
2. **Market Data Loading**: Fetches historical Bitcoin data from Polygon.io
3. **Phase Transition**: Automatically transitions to Lobby phase
4. **Player Connection**: Players connect via WebSocket and receive game state

#### Live Trading Flow

1. **Player Action**: Player creates/closes position via REST API
2. **Service Processing**: PlayerService validates and processes the action
3. **State Update**: Player session is updated with new position/balance
4. **Real-time Broadcast**: Updates are sent to all connected players
5. **P&L Calculation**: Real-time P&L updates sent to affected players

#### Price Update Flow

1. **Market Data**: RoundManager streams historical data as live prices
2. **Price Broadcast**: New price data sent to all connected clients
3. **P&L Update**: PlayerService recalculates P&L for all active positions
4. **Individual Updates**: P&L updates sent to each affected player
5. **Leaderboard Update**: Leaderboard recalculated and broadcast to all players

## Frontend Architecture

### State Management

The frontend uses Zustand for state management with three main stores:

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend State                          │
├─────────────────────────────────────────────────────────────┤
│  useAuthStore  │  useGameStore  │  useWsStore             │
│                │                │                          │
│ • JWT Tokens   │ • Game State   │ • WebSocket Connection  │
│ • Login/Logout │ • Positions    │ • Connection Status     │
│ • Token Refresh│ • Chart Data   │ • Message Routing       │
│ • User Data    │ • P&L Data     │ • Reconnection Logic    │
│                │ • Leaderboard  │                          │
└─────────────────────────────────────────────────────────────┘
```

### Component Architecture

```
App
├── Layout
│   ├── JoinGameOverlay (Modal)
│   └── Main Game Interface
│       ├── CandlestickChart (Real-time chart)
│       ├── Leaderboard (Player rankings)
│       ├── TradingPanel (Position controls)
│       ├── GameInfo (Phase & stats)
│       └── Timer (Countdown)
```

### Real-time Communication

#### WebSocket Message Flow

1. **Connection**: Frontend establishes WebSocket connection with JWT token
2. **Initial Sync**: Receives complete game state on connection
3. **Message Handling**: Each message type updates specific store state
4. **UI Updates**: React components re-render based on state changes

#### API Integration

1. **Authentication**: Login endpoint returns JWT tokens
2. **Position Management**: REST API calls for position creation/closing
3. **Error Handling**: Comprehensive error handling with user feedback
4. **Token Refresh**: Automatic token refresh on expiration

## Data Models

### Core Domain Models

#### Player State

```go
type PlayerState struct {
    PlayerId string
    BasePlayerState
}

type BasePlayerState struct {
    Balance         float64
    ActivePosition  *Position
    ClosedPositions []ClosedPosition
}
```

#### Position

```go
type Position struct {
    Type          PositionType
    EntryPrice    float64
    EntryTime     time.Time
    Quantity      float64
    Pnl           float64
    PnlPercentage float64
}
```

#### Game State

```go
type GameStatePayload struct {
    RoundID             string
    ChartData           []PriceData
    Phase               Phase
    EndTime             time.Time
    Balance             float64
    ActivePosition      *Position
    ClosedPositions     []ClosedPosition
    TotalPnl            float64
    ActivePnl           float64
    ActivePnlPercentage float64
    LongPositions       int
    ShortPositions      int
    TotalPlayers        int
}
```

## Concurrency & Thread Safety

### Backend Concurrency Model

The backend uses Go's goroutines and channels for concurrent operations:

#### Thread-Safe Operations

- **PlayerService**: Uses read/write locks for session management
- **RoundManager**: Uses read/write locks for game state
- **Hub**: Uses channels for message passing between goroutines

#### Concurrent Operations

- **WebSocket Connections**: Each client runs in separate goroutines
- **Market Data Streaming**: Price updates run in background goroutine
- **P&L Calculations**: Real-time updates run concurrently
- **Message Broadcasting**: Non-blocking message distribution

### Frontend Concurrency

- **WebSocket Connection**: Single connection per player
- **State Updates**: Zustand handles concurrent state updates
- **API Calls**: Async/await pattern for REST API calls
- **Real-time Updates**: Non-blocking message processing

## Scalability Considerations

### Current Implementation

- **In-Memory Sessions**: Player sessions are stored in memory for fast access
- **Concurrent Operations**: Thread-safe operations using read/write locks
- **WebSocket Optimization**: Efficient message broadcasting to all connected clients
- **Database Minimalism**: Database is used primarily for player authentication, not game state

### Current Limitations

1. **In-Memory Sessions**: Player sessions stored in memory
2. **Single Server**: No horizontal scaling capability
3. **Database Bottleneck**: Could become limiting factor
4. **WebSocket Limits**: No connection pooling or load balancing

### Future Scalability Plans

1. **Redis Sessions**: Move sessions to Redis for multi-server support
2. **Microservices**: Split into separate services (auth, trading, market data)
3. **Load Balancing**: Add load balancer for WebSocket connections
4. **Database Optimization**: Add caching and connection pooling

## Security Architecture

### Authentication & Authorization

#### JWT Implementation

- **Access Token**: Short-lived (15 minutes) for API requests
- **Refresh Token**: Long-lived (7 days) for token refresh
- **Token Storage**: Frontend stores in memory, backend validates on each request

#### Security Measures

- **Input Validation**: All inputs validated on server side
- **SQL Injection Protection**: Parameterized queries
- **XSS Protection**: Frontend escaping and sanitization
- **CORS Configuration**: Proper CORS headers for cross-origin requests

### Data Protection

#### Sensitive Data

- **Player Credentials**: Not stored (username-based authentication)
- **Trading Data**: In-memory only, not persisted
- **Market Data**: Public data from Polygon.io API

#### Privacy Considerations

- **No Personal Data**: Only usernames are collected
- **Session Data**: Cleared on server restart
- **Trading History**: Not persisted between rounds

## Deployment Architecture

### Docker Containerization

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Compose                          │
├─────────────────────────────────────────────────────────────┤
│  Frontend Container  │  Backend Container  │  DB Container │
│  (Next.js)          │  (Go)               │  (PostgreSQL) │
│  Port: 3000         │  Port: 8080         │  Port: 5432   │
└─────────────────────────────────────────────────────────────┘
```

### Environment Configuration

#### Backend Environment

```env
PORT=8080
DATABASE_URL=postgres://user:password@postgres:5432/tradeoff
POLYGON_API_KEY=your_api_key
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION=15m
```

#### Frontend Environment

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Monitoring & Observability

### Current Monitoring

#### Backend Metrics

- **Player Count**: Real-time active player tracking
- **Position Counts**: Long/short position statistics
- **Game Phases**: Phase transition logging
- **Error Logging**: Comprehensive error logging

#### Frontend Metrics

- **Connection Status**: WebSocket connection health
- **API Response Times**: REST API performance
- **User Interactions**: Trading action tracking
- **Error Tracking**: Client-side error reporting

### Future Monitoring Plans

#### Backend Enhancements

- **Prometheus Metrics**: Custom metrics for game performance
- **Distributed Tracing**: Request tracing across services
- **Health Checks**: Comprehensive health check endpoints
- **Performance Profiling**: CPU and memory profiling

#### Frontend Enhancements

- **Analytics Integration**: User behavior tracking
- **Performance Monitoring**: Core Web Vitals tracking
- **Error Reporting**: Automated error reporting
- **A/B Testing**: Feature flag implementation

## Development Workflow

### Code Organization

#### Backend Structure

```
backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── domain/         # Core data models
│   ├── handler/        # HTTP/WebSocket handlers
│   ├── middleware/     # HTTP middleware
│   ├── platform/       # Infrastructure (router)
│   ├── service/        # Business logic
│   └── storage/        # Data persistence
├── config/             # Configuration files
└── Dockerfile          # Container definition
```

#### Frontend Structure

```
frontend/
├── src/
│   ├── app/            # Next.js app router
│   ├── components/     # React components
│   ├── stores/         # Zustand state management
│   └── types.d.ts      # TypeScript definitions
├── public/             # Static assets
└── Dockerfile          # Container definition
```

### Development Practices

#### Backend Practices

- **Clean Architecture**: Clear separation of concerns
- **Error Handling**: Comprehensive error handling
- **Logging**: Structured logging for debugging
- **Testing**: Unit tests for core business logic

#### Frontend Practices

- **Component Design**: Reusable, composable components
- **State Management**: Centralized state with Zustand
- **Type Safety**: Comprehensive TypeScript usage
- **Performance**: Optimized rendering and updates

## Future Architecture Evolution

### Phase 2 Enhancements

1. **Live Leaderboard**: Real-time display of active user count and player rankings ✅ **COMPLETE**
2. **Player Rankings**: Live leaderboard showing top performers ✅ **COMPLETE**
3. **Performance Metrics**: Trading statistics and rankings ✅ **COMPLETE**

### Phase 3 Microservices

1. **Auth Service**: Dedicated authentication service
2. **Trading Service**: Position management microservice
3. **Market Data Service**: Real-time market data service
4. **Game Service**: Game state management service
5. **Leaderboard Service**: Dedicated leaderboard and ranking service
6. **Notification Service**: Real-time notifications
7. **Analytics Service**: Trading analytics and reporting

### Infrastructure Evolution

1. **Load Balancing**: Multiple backend instances
2. **Caching**: Redis for session and data caching
3. **Message Queues**: NATS/RabbitMQ for async processing
4. **Monitoring**: Prometheus, Grafana, distributed tracing
5. **CI/CD**: Automated testing and deployment pipelines

## Conclusion

The current TradeOff architecture provides a solid foundation for a multiplayer trading game with real-time capabilities. The monolithic design was chosen for Phase 1 to enable rapid development and iteration. The system demonstrates good separation of concerns, thread safety, and real-time communication capabilities.

As the application scales, the architecture is designed to evolve into a microservices-based system with proper load balancing, caching, and monitoring. The current implementation provides valuable insights into performance characteristics and scalability requirements for future phases.
