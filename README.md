# TradeOff: The 10-Minute Stock Market Game

![Go Badge](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Next.js Badge](https://img.shields.io/badge/Next.js-000000?style=for-the-badge&logo=next.js&logoColor=white)
![React Badge](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![Docker Badge](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![PostgreSQL Badge](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)

Welcome to **TradeOff**! This is a real-time, massively multiplayer stock market simulator designed to provide all the thrill of high-stakes trading in compressed, 10-minute rounds, with zero financial risk.

This project is being built in public as a comprehensive case study in full-stack application development, product design, and scalable system architecture. You can follow the progress and my learnings on [Twitter](https://twitter.com/jassiD2000) and [LinkedIn](https://www.linkedin.com/in/jaswinder-singh-32a01118b/).

**Live Demo:** `[Link to your deployed frontend - Coming Soon!]`

---

## Current Project Status (Phase 2 Complete)

**Phase 2 has been successfully completed!** The game now features a live leaderboard that displays real-time player rankings and active user statistics. Players can see their position on the leaderboard and track their performance against other players in real-time.

### 🎯 Phase 2 Achievements

**Live Leaderboard System Implemented:**

- ✅ **Real-time Rankings**: Live leaderboard showing top 20 players sorted by active balance
- ✅ **Player Statistics**: Real-time display of total players online and position distributions
- ✅ **Personal Highlighting**: Current player's position is highlighted on the leaderboard
- ✅ **Live Updates**: Leaderboard updates automatically as players' balances change
- ✅ **Balance Calculation**: Active balance includes unrealized P&L from open positions
- ✅ **WebSocket Integration**: New `leaderboard_update` message type for real-time updates
- ✅ **Backend Service**: `GetLeaderboard()` method in PlayerService for ranking calculations
- ✅ **Frontend Component**: New `Leaderboard.tsx` component with responsive design
- ✅ **UI Integration**: Leaderboard prominently displayed alongside the trading chart

**Technical Implementation:**

- **Backend**: Enhanced `PlayerService` with leaderboard calculation logic
- **Frontend**: New leaderboard component with real-time updates via WebSocket
- **Data Flow**: Leaderboard updates triggered during P&L calculations
- **Performance**: Limited to top 20 players for optimal performance
- **User Experience**: Seamless integration with existing game interface

### ✅ Completed Features

#### Core Game Mechanics

- **Multiplayer Game Sessions**: Multiple players can join and participate in the same trading round simultaneously
- **Real-time Position Management**: Players can create long/short positions during live trading phases
- **Dynamic P&L Tracking**: Real-time profit/loss calculation based on current market prices
- **Position Controls**: Intuitive UI for position creation, management, and closing
- **Round-based Gameplay**: Continuous 3-phase game loop (Lobby → Live → Cooldown)

#### Live Leaderboard System

- **Real-time Rankings**: Live leaderboard showing top 20 players sorted by active balance
- **Player Statistics**: Real-time display of total players online and position distributions
- **Personal Highlighting**: Current player's position is highlighted on the leaderboard
- **Live Updates**: Leaderboard updates automatically as players' balances change
- **Balance Calculation**: Active balance includes unrealized P&L from open positions

#### Real-time Communication

- **WebSocket Integration**: Real-time game state updates, price data, and position changes
- **Live Chart Updates**: Interactive candlestick charts that update in real-time
- **Phase Transitions**: Seamless transitions between game phases with countdown timers
- **Player Statistics**: Real-time tracking of player counts, position distributions, and balances

#### Technical Infrastructure

- **JWT Authentication**: Secure player authentication and session management
- **Market Data Integration**: Real Bitcoin/USD price data from Polygon.io API
- **Containerized Deployment**: Full Docker support for consistent development and deployment
- **Database Persistence**: PostgreSQL for player data and session management

#### User Experience

- **Responsive Design**: Modern, mobile-friendly interface built with Tailwind CSS
- **Interactive Trading Panel**: Clear position controls with real-time P&L display
- **Game Information Display**: Live phase status, player counts, and position statistics
- **Seamless Onboarding**: Simple username-based player registration
- **Leaderboard Integration**: Prominent leaderboard display alongside the trading chart

---

## The Vision: A Three-Phase Architectural Journey

This project is planned to evolve through three distinct phases, demonstrating a professional approach to building and scaling a modern web application.

### Phase 1: The Multiplayer Core ✅ **COMPLETE**

- **Goal:** Build a functional, engaging, multiplayer version of the game with real-time trading capabilities.
- **Architecture:** Monolith (Single Go Backend + Next.js Frontend).
- **Status:** ✅ **Complete** - Full multiplayer support with real-time trading

### Phase 2: The Live Leaderboard ✅ **COMPLETE**

- **Goal:** Add a live leaderboard that displays active user count and player rankings in real-time.
- **Architecture:** The application remains a **Monolith**, enhanced with leaderboard functionality.
- **Status:** ✅ **Complete** - Live leaderboard with real-time player rankings

### Phase 3: The Scalability Case Study ⏳ **NEXT**

- **Goal:** Use professional load-testing tools (`k6`) to prove the monolith's performance limits under simulated multiplayer load, and then evolve the architecture to solve those bottlenecks.
- **Architecture:** Refactor the system into an **Evolved Microservices Architecture**. Key components (e.g., data ingestion, trade processing) will be extracted into separate, scalable services.
- **Status:** ⏳ **Planned**

_[This README will be updated with "Before vs. After" performance benchmark graphs upon completion of Phase 3.]_

---

## Tech Stack Overview

| Component          | Technology                               |
| :----------------- | :--------------------------------------- |
| **Frontend**       | Next.js, React, TypeScript, Tailwind CSS |
| **Backend**        | Go, Chi (Router), Gorilla WebSocket      |
| **Database**       | PostgreSQL                               |
| **Authentication** | JWT (JSON Web Tokens)                    |
| **Market Data**    | Polygon.io API (Bitcoin/USD)             |
| **Orchestration**  | Docker, Docker Compose                   |
| **Planned**        | Redis (Phase 2), NATS/RabbitMQ (Phase 3) |

---

## Getting Started

This project is fully containerized with Docker, providing a consistent and easy-to-manage development environment.

### Prerequisites

- [Git](https://git-scm.com/downloads)
- [Docker](https://www.docker.com/products/docker-desktop/) & Docker Compose

### Local Development Setup

1.  **Clone the Repository:**

    ```bash
    git clone https://github.com/jassi-singh/tradeoff.git
    cd tradeoff
    ```

2.  **Configure Backend Environment:**
    The backend requires a `.env` file with credentials for the database and the Polygon.io API.

    First, copy the example file:

    ```bash
    cp backend/.env.example backend/.env
    ```

    Next, **edit `backend/.env`** and add your Polygon API key. The `DATABASE_URL` is already configured for the Docker environment.

    ```env
    # backend/.env

    # Port for the backend server
    PORT=8080

    # PostgreSQL Connection URL (configured for Docker)
    DATABASE_URL="postgres://user:password@postgres:5432/tradeoff?sslmode=disable"

    # API Key for Polygon.io (REQUIRED)
    POLYGON_API_KEY="YOUR_POLYGON_API_KEY" # <--- ADD YOUR KEY HERE
    ```

3.  **Configure Frontend Environment:**
    The frontend needs to know the backend's URL. Copy the example file. No edits are needed for the default Docker setup.

    ```bash
    cp frontend/.env.example frontend/.env.local
    ```

4.  **Run the Application:**
    Use Docker Compose to build the images and start all services (frontend, backend, and database).

    ```bash
    docker-compose up --build
    ```

    - The frontend will be available at `http://localhost:3000`.
    - The backend API will be available at `http://localhost:8080`.
    - The PostgreSQL database will be accessible on port `5432`.

    Once the containers are running, the backend will automatically create the necessary database tables.

### Alternative: Running Services Manually

If you prefer not to use Docker, you can run each service manually. For detailed instructions, see the README file within its directory:

- **[Backend README](./backend/README.md)**
- **[Frontend README](./frontend/README.md)**

---

## Game Rules & Mechanics

### Game Phases

1. **Lobby (15 seconds)**: Players join and wait for the trading round to begin
2. **Live (1 minute)**: Active trading phase where players can create and close positions
3. **Cooldown (10 seconds)**: Brief pause before the next round begins

### Trading Mechanics

- **Starting Balance**: $100 USD per player per round
- **Position Types**: Long (betting price goes up) or Short (betting price goes down)
- **Position Limits**: One active position per player at a time
- **P&L Calculation**: Real-time based on current market price vs entry price
- **Market Data**: Real Bitcoin/USD historical data from Polygon.io

### Multiplayer Features

- **Concurrent Players**: Multiple players can join the same game session
- **Real-time Updates**: All players see live price movements and position changes
- **Player Statistics**: Live tracking of total players, long/short position counts
- **Session Management**: Players maintain their positions and balances throughout the round

---

## API Documentation

For comprehensive API documentation, including all endpoints, WebSocket messages, and data structures, see **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)**.

---

## Contributing

This project is built in public as a learning experience. Feel free to:

- **Follow the Journey**: Check out the [Twitter](https://twitter.com/jassiD2000) and [LinkedIn](https://www.linkedin.com/in/jaswinder-singh-32a01118b/) for updates
- **Try the Game**: Run it locally and experience the multiplayer trading
- **Learn from the Code**: Study the architecture and implementation patterns
- **Provide Feedback**: Share thoughts on the design decisions and implementation

---

## License

This project is open source and available under the [MIT License](LICENSE).
