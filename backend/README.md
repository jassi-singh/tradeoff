# TradeOff - Backend Service

![Go Badge](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Chi Badge](https://img.shields.io/badge/chi-800080?style=for-the-badge)
![WebSocket Badge](https://img.shields.io/badge/WebSocket-019537?style=for-the-badge&logo=websocket&logoColor=white)
![PostgreSQL Badge](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker Badge](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)

This repository contains the backend service for **TradeOff**, a real-time, massively multiplayer stock market simulator.

This service is written in Go and is responsible for managing player state, handling API requests, and serving real-time game data via WebSockets to the frontend client.

---

## Core Features

The backend implements a complete, multiplayer, round-based trading game cycle.

### Game Mechanics

- **Multiplayer Game Sessions**: Multiple players can join and participate in the same trading round simultaneously
- **Real-time Game Phases**: The game operates in a continuous loop of three phases:
  1.  **Lobby (15s):** A waiting period where players join and the system loads historical market data for the upcoming round.
  2.  **Live (60s):** The active trading phase. Players can create and close positions while the backend streams historical price data to all connected clients in real-time, simulating a live market.
  3.  **Cooldown (10s):** A brief period after the live phase before a new round begins in the lobby.
- **Position Management**: Players can create long/short positions and close them during live trading phases
- **Real-time P&L Tracking**: Live profit/loss calculation based on current market prices for all players

### Real-time Communication

- **WebSocket Communication**: Real-time game state, chart data, price updates, and position changes are pushed to clients using WebSockets
- **Player Authentication**: JWT-based authentication system for secure player sessions
- **Concurrent Player Support**: Handles multiple simultaneous WebSocket connections with proper session management

### Market Data & Infrastructure

- **Market Data Simulation**: The service uses the [Polygon.io](https://polygon.io/) API to fetch real historical price data for `X:BTCUSD` (Bitcoin/USD), which is then used to simulate the game's market
- **Player Management**: A REST API is available to create and retrieve players, with data persisted in a PostgreSQL database
- **Session Management**: In-memory player sessions with concurrent-safe operations for real-time game state

---

## Technology Stack

- **Language:** Go
- **Web Router:** [Chi](https://github.com/go-chi/chi)
- **WebSockets:** [Gorilla WebSocket](https://github.com/gorilla/websocket)
- **Database:** PostgreSQL
- **Configuration:** [Viper](https://github.com/spf13/viper)
- **Authentication:** JWT (JSON Web Tokens)
- **Market Data:** [Polygon.io API](https://polygon.io/)

---

## Project Architecture

The project follows a clean, layered architecture to separate concerns and improve maintainability:

- `/cmd/server`: The main application entry point. Responsible for initializing dependencies (database, services) and starting the HTTP server.
- `/internal/config`: Handles loading application configuration from `config.yml` and environment variables.
- `/internal/domain`: Contains the core data structures (models) of the application, such as `Player`, `Position`, `PriceData`, and game phases.
- `/internal/handler`: The web layer. Contains HTTP and WebSocket handlers responsible for processing incoming requests and interacting with the service layer.
  - `auth_handler.go`: Handles player authentication and JWT token management
  - `position_handler.go`: Manages position creation and closing operations
  - `websocket_handler.go`: Handles WebSocket connections and real-time communication
- `/internal/service`: Contains the core business logic.
  - `round_manager.go`: Manages the game state, phase transitions, and the main game loop
  - `market_service.go`: Fetches data from the Polygon.io API
  - `player_service.go`: Manages player sessions, positions, and P&L calculations
  - `hub.go`: Manages all active WebSocket client connections
  - `auth_service.go`: Handles JWT token generation and validation
- `/internal/platform/router`: Configures the Chi router and defines all API routes.
- `/internal/storage`: The data persistence layer. Implements repository interfaces for interacting with the PostgreSQL database.
- `/internal/middleware`: HTTP middleware for authentication and request processing.

---

## API Documentation

For comprehensive API documentation, including all endpoints, WebSocket messages, and data structures, see **[API_DOCUMENTATION.md](../API_DOCUMENTATION.md)**.

---

## Getting Started

Follow these instructions to get the backend server running on your local machine.

### Prerequisites

- [Go](https://go.dev/doc/install) (version 1.22 or newer)
- [Docker](https://www.docker.com/get-started/) (for running a local PostgreSQL instance)
- A code editor like [VSCode](https://code.visualstudio.com/)

### Installation & Setup

1.  **Clone the Repository:**
    This backend is part of a larger monorepo. Clone the main project to get started.

    ```bash
    git clone https://github.com/jassi-singh/tradeoff.git
    cd tradeoff/backend
    ```

2.  **Set Up Environment Variables:**
    Create a `.env` file in the `/backend` directory by copying the example below.

    ```bash
    cp .env.example .env
    ```

    You will need to create a file named `.env.example` with the following content:

    ```
    # Port for the backend server
    PORT=8080

    # PostgreSQL Connection URL
    DATABASE_URL="postgres://user:password@localhost:5432/tradeoff?sslmode=disable"

    # API Key for Polygon.io
    POLYGON_API_KEY="YOUR_POLYGON_API_KEY"
    ```

    Update the `DATABASE_URL` with your database credentials and add your [Polygon.io API key](https://polygon.io/dashboard).

3.  **Set Up Database:**
    You can run a PostgreSQL instance using Docker:

    ```bash
    docker run --name tradeoff-db -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=tradeoff -p 5432:5432 -d postgres
    ```

    Once the database is running, connect to it and create the `players` table:

    ```sql
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    CREATE TABLE players (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        username VARCHAR(255) NOT NULL,
        refresh_token TEXT,
        refresh_token_expiry TIMESTAMP
    );
    ```

4.  **Install Dependencies:**
    This command will download the packages listed in your `go.mod` file.
    ```bash
    go mod tidy
    ```

### Running the Server

#### With Go

To start the server directly with Go, run the following command from the `/backend` directory:

```bash
go run ./cmd/server/main.go
```

The server will start on the port specified in your `.env` file (default: 8080).

#### With Docker

You can also build and run the backend service using the provided Dockerfile.

1.  **Build the Docker Image:**
    From the `tradeoff/backend` directory, run:

    ```bash
    docker build -t tradeoff-backend .
    ```

2.  **Run the Docker Container:**
    Make sure to pass your local `.env` file to the container.
    ```bash
    docker run --name tradeoff-backend --env-file .env -p 8080:8080 -d tradeoff-backend
    ```
    **Note:** If your PostgreSQL database is also running in a Docker container, you may need to adjust the `DATABASE_URL` in your `.env` file to point to the database container's network address (e.g., using `host.docker.internal` or a shared Docker network).

---

## Game Logic Details

### Round Management

- **Round Duration**: 85 seconds total (15s lobby + 60s live + 10s cooldown)
- **Continuous Loop**: Rounds automatically restart after cooldown phase
- **Market Data**: Each round uses different historical Bitcoin data for variety

### Player Sessions

- **Session Creation**: Players get a session when they first connect via WebSocket
- **Concurrent Safety**: All session operations are thread-safe using read/write locks
- **Session Reset**: Player sessions are reset at the start of each new round

### Position Management

- **One Position Per Player**: Players can only have one active position at a time
- **Full Balance Investment**: When creating a position, the player invests their entire balance
- **Real-time P&L**: P&L is calculated and updated in real-time during live trading
- **Position Types**: Long (profit when price goes up) and Short (profit when price goes down)

### WebSocket Communication

- **Connection Management**: Hub manages all active WebSocket connections
- **Message Broadcasting**: Game state updates are broadcast to all connected players
- **Direct Messages**: P&L updates are sent directly to individual players
- **Connection Cleanup**: Proper cleanup when players disconnect

---

## Development Notes

### Key Design Decisions

1. **Monolithic Architecture**: Chosen for Phase 1 simplicity and rapid development
2. **In-Memory Game State**: Provides fast access and real-time updates
3. **JWT Authentication**: Stateless authentication suitable for WebSocket connections
4. **Real Market Data**: Uses actual Bitcoin price data for authentic trading experience

### Testing Strategy

- **Unit Tests**: Core business logic in services
- **Integration Tests**: API endpoints and WebSocket communication
- **Load Testing**: Planned for Phase 3 to identify bottlenecks

### Monitoring & Logging

- **Structured Logging**: All major operations are logged with appropriate detail
- **Error Handling**: Comprehensive error handling with meaningful error messages
- **Game Metrics**: Basic metrics for player counts and position statistics

---

## Future Enhancements

### Phase 2: Live Leaderboard

- **Active User Count**: Real-time display of connected players
- **Player Rankings**: Live leaderboard showing top performers
- **Performance Metrics**: Trading statistics and rankings

### Phase 3: Scalability Study

- **Load Testing**: Performance analysis under high load
- **Microservices**: Architecture evolution based on bottlenecks
- **Advanced Monitoring**: Comprehensive performance monitoring
