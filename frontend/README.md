# TradeOff - Frontend Service

![Next.js Badge](https://img.shields.io/badge/Next.js-000000?style=for-the-badge&logo=next.js&logoColor=white)
![React Badge](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![Tailwind CSS Badge](https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)
![Zustand Badge](https://img.shields.io/badge/Zustand-4B3A29?style=for-the-badge)
![Docker Badge](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)

This is the frontend client for **TradeOff**, a real-time, massively multiplayer stock market simulator.

This application is built with Next.js and React. It is responsible for rendering the game interface, handling user interactions, and communicating with the Go backend service via REST APIs and WebSockets.

---

## Core Features

### Multiplayer Game Experience

- **Real-time Gameplay:** Connects to the backend via WebSockets to receive live game state updates, including phase transitions (Lobby, Live, Cooldown) and market data
- **Multiplayer Sessions:** Multiple players can join the same game and see each other's activities in real-time
- **Player Onboarding:** A simple and clean overlay allows new players to join a game by entering a username. Player sessions are persisted in local storage
- **Live Player Statistics:** Real-time display of total players, long/short position counts, and game phase information

### Trading Interface

- **Interactive Charting:** Displays historical and real-time stock data using **Lightweight Charts**, which updates live as new price data is streamed from the backend
- **Position Management:** Intuitive UI for creating long/short positions and closing them during live trading phases
- **Real-time P&L Tracking:** Live profit/loss calculation and display with color-coded indicators
- **Trading Panel:** Clear position controls with real-time balance and P&L information

### User Experience

- **Dynamic UI:** The UI dynamically reflects the current game state, showing the active phase and a countdown timer for the next phase transition
- **Responsive Design:** Modern, mobile-friendly interface built with Tailwind CSS
- **Component-Based Architecture:** Built with reusable React components for different parts of the UI like the chart, timer, trading panel, and game info display
- **State Management:** Comprehensive state management using Zustand for game state, authentication, and WebSocket connections

---

## Technology Stack

- **Framework:** [Next.js](https://nextjs.org/) (with App Router)
- **Library:** [React](https://react.dev/)
- **Styling:** [Tailwind CSS](https://tailwindcss.com/)
- **State Management:** [Zustand](https://github.com/pmndrs/zustand)
- **Charting Library:** [Lightweight Charts](https://www.tradingview.com/lightweight-charts/)
- **Authentication:** JWT token management with automatic refresh
- **Real-time Communication:** WebSocket integration for live updates

---

## Project Architecture

The project's `src` directory is organized to separate concerns:

- `/app`: Contains the main page and layout definitions for the Next.js App Router
- `/api.ts`: A dedicated file for functions that interact with the backend REST and WebSocket APIs
- `/components`: Contains all reusable React components:
  - `CandlestickChart.tsx`: Interactive real-time chart component
  - `TradingPanel.tsx`: Position management and trading controls
  - `GameInfo.tsx`: Game phase and player statistics display
  - `Timer.tsx`: Countdown timer for phase transitions
  - `overlay/`: Modal components for player onboarding
- `/stores`: Zustand state management stores:
  - `useAuthStore.ts`: Authentication state and token management
  - `useGameStore.ts`: Game state, positions, and real-time updates
  - `useWsStore.ts`: WebSocket connection management
- `/types.d.ts`: TypeScript type definitions for shared data structures

---

## Getting Started

Follow these instructions to get the frontend development server running on your local machine.

### Prerequisites

- [Node.js](https://nodejs.org/en) (version 20 or newer)
- [pnpm](https://pnpm.io/installation) (recommended package manager)

### Installation & Setup

1.  **Clone the Repository:**
    This frontend lives in the `/frontend` directory of the main monorepo.

    ```bash
    git clone https://github.com/jassi-singh/tradeoff.git
    cd tradeoff/frontend
    ```

2.  **Set Up Environment Variables:**
    Create a `.env.local` file by copying the example. This will be used to store the URL of your backend server.

    ```bash
    cp .env.example .env.local
    ```

    You will need to create a file named `.env.example` with the following content:

    ```
    # The base URL for the backend API and WebSocket server
    # Default is for a local backend running on port 8080
    NEXT_PUBLIC_API_URL="http://localhost:8080"
    ```

3.  **Install Dependencies:**
    This command will install all the necessary packages from `package.json`.
    ```bash
    pnpm install
    ```

### Running the Development Server

To start the local development server, run the following command:

```bash
pnpm dev
```

The application will be available at `http://localhost:3000`.

---

## Game Interface Components

### Main Game Layout

The game interface is divided into several key components:

1. **Chart Area**: Real-time candlestick chart showing Bitcoin/USD price movements
2. **Trading Panel**: Position management controls and P&L display
3. **Game Information**: Current phase, timer, and player statistics
4. **Player Onboarding**: Modal overlay for new player registration

### Real-time Features

- **Live Price Updates**: Chart updates in real-time as new price data arrives
- **Phase Transitions**: UI automatically updates when game phases change
- **Position Updates**: Trading panel reflects current position status and P&L
- **Player Counts**: Live updates of total players and position distributions

### Trading Interface

- **Position Creation**: "Go Long" and "Go Short" buttons during live phase
- **Position Management**: Active position display with entry price, time, and current P&L
- **Position Closing**: "Close Position" button to exit current position
- **Balance Tracking**: Real-time balance and P&L updates

---

## State Management

### Authentication Store (`useAuthStore`)

- Manages JWT tokens (access and refresh)
- Handles login/logout functionality
- Automatic token refresh on expiration
- Persistent authentication state

### Game Store (`useGameStore`)

- Game state (phase, round ID, end time)
- Player data (balance, positions, P&L)
- Chart data and real-time updates
- Trading actions (create/close positions)

### WebSocket Store (`useWsStore`)

- WebSocket connection management
- Connection status tracking
- Automatic reconnection handling
- Message routing to game store

---

## API Integration

### REST API Calls

- **Authentication**: Login and token refresh endpoints
- **Position Management**: Create and close position endpoints
- **Error Handling**: Comprehensive error handling with user feedback

### WebSocket Communication

- **Connection**: Automatic WebSocket connection with JWT authentication
- **Message Handling**: Real-time message processing for game updates
- **Reconnection**: Automatic reconnection on connection loss
- **Message Types**: Handles all game state, price, and P&L updates

---

## User Experience Features

### Responsive Design

- **Mobile-First**: Optimized for mobile and desktop viewing
- **Dark Theme**: Modern dark theme for trading interface
- **Accessibility**: Proper contrast and keyboard navigation support

### Real-time Feedback

- **Loading States**: Visual feedback during API calls
- **Error Handling**: User-friendly error messages
- **Success Indicators**: Clear feedback for successful actions

### Game Flow

1. **Player Registration**: Simple username entry to join game
2. **Game Connection**: Automatic WebSocket connection and game state sync
3. **Trading Phase**: Active trading during live phase with position controls
4. **Round Transitions**: Seamless transitions between game phases

---

## Development Notes

### Key Design Decisions

1. **Component Architecture**: Modular components for maintainability
2. **State Management**: Zustand for simple, efficient state management
3. **Real-time Updates**: WebSocket-first design for live game experience
4. **Type Safety**: Comprehensive TypeScript types for all data structures

### Development Practices

- **Component Design**: Reusable, composable components
- **State Management**: Centralized state with Zustand
- **Type Safety**: Comprehensive TypeScript usage
- **Performance**: Optimized rendering and updates

### Testing Strategy

- **Component Testing**: Unit tests for individual components
- **Integration Testing**: End-to-end testing of game flow
- **State Testing**: Testing of Zustand store logic
- **API Testing**: Testing of API integration and error handling

---

## Deployment

### Docker Support

The frontend includes a Dockerfile for containerized deployment:

```bash
docker build -t tradeoff-frontend .
docker run -p 3000:3000 tradeoff-frontend
```

### Environment Configuration

- **Development**: Uses local backend URL
- **Production**: Configurable backend URL via environment variables
- **Docker**: Pre-configured for Docker Compose deployment

### Build Optimization

- **Next.js Optimization**: Automatic code splitting and optimization
- **Bundle Analysis**: Tools for analyzing bundle size and performance
- **Image Optimization**: Automatic image optimization for better performance

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
