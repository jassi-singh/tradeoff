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

- **Real-time Gameplay:** Connects to the backend via WebSockets to receive live game state updates, including phase transitions (Lobby, Live, Cooldown) and market data.
- **Player Onboarding:** A simple and clean overlay allows new players to join a game by entering a username. Player sessions are persisted in local storage.
- **Interactive Charting:** Displays historical and real-time stock data using **Lightweight Charts**, which updates live as new price data is streamed from the backend.
- **Dynamic UI:** The UI dynamically reflects the current game state, showing the active phase and a countdown timer for the next phase transition.
- **Component-Based Architecture:** Built with reusable React components for different parts of the UI like the leaderboard, portfolio, and game info display.

---

## Technology Stack

- **Framework:** [Next.js](https://nextjs.org/) (with App Router)
- **Library:** [React](https://react.dev/)
- **Styling:** [Tailwind CSS](https://tailwindcss.com/)
- **State Management:** [Zustand](https://github.com/pmndrs/zustand)
- **Charting Library:** [Lightweight Charts](https://www.tradingview.com/lightweight-charts/)

---

## Project Architecture

The project's `src` directory is organized to separate concerns:

- `/app`: Contains the main page and layout definitions for the Next.js App Router.
- `/api.ts`: A dedicated file for functions that interact with the backend REST and WebSocket APIs.
- `/components`: Contains all reusable React components, such as the chart, timer, and UI overlays.
- `/hooks`: Custom React hooks, like `useWS` for simplifying WebSocket message handling.
- `/stores`: Zustand state management stores, which handle global client state for authentication (`useAuthStore`), WebSocket connections (`useWsStore`), and game state (`useGameStore`).
- `/types.d.ts`: TypeScript type definitions for shared data structures.

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
