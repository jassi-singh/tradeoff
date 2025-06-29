# TradeOff - Backend Service

![Go Badge](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Chi Badge](https://img.shields.io/badge/chi-800080?style=for-the-badge)
![WebSocket Badge](https://img.shields.io/badge/WebSocket-019537?style=for-the-badge&logo=websocket&logoColor=white)
![PostgreSQL Badge](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)

This repository contains the backend service for **TradeOff**, a real-time, massively multiplayer stock market simulator.

This service is written in Go and is responsible for managing player state, handling API requests, and serving real-time game data via WebSockets to the frontend client.

---

## Project Status: Phase 1 (Initial Scaffolding)

The project is in its initial setup phase. The directory structure has been established according to clean architecture principles to ensure maintainability and future scalability. The current codebase provides the foundational boilerplate for the HTTP server, API routing, and WebSocket connections.

The immediate goal is to implement the core logic within these foundational files.

---

## Technology Stack

* **Language:** Go
* **Web Router:** [Chi](https://github.com/go-chi/chi)
* **WebSockets:** [Gorilla WebSocket](https://github.com/gorilla/websocket)
* **Database (Planned):** PostgreSQL

---

## Directory Structure Overview

The project follows a clean, layered architecture:

* `/cmd/server`: The main application entry point. Responsible for initializing dependencies and starting the server.
* `/internal/domain`: Contains the core data structures (models) of the application.
* `/internal/game`: Will contain the core business logic and state management for the game itself.
* `/internal/handler`: The web layer. Contains all HTTP and WebSocket handlers responsible for handling incoming requests.
* `/internal/platform/router`: Configures the Chi router and defines all API routes.
* `/internal/storage`: (To be added) Will handle all database interactions.

---

## Getting Started

Follow these instructions to get the backend server running on your local machine.

### Prerequisites

* [Go](https://go.dev/doc/install) (version 1.22 or newer)
* A code editor like [VSCode](https://code.visualstudio.com/)

### Installation & Setup

1.  **Clone the Repository:**
    This backend is part of a larger monorepo. Clone the main project to get started.
    ```bash
    git clone [https://github.com/jassi-singh/tradeoff.git](https://github.com/jassi-singh/tradeoff.git)
    cd tradeoff/backend
    ```

2.  **Set Up Environment Variables:**
    Copy the example environment file. The server will run on port 8080 by default if the `.env` file is not present.
    ```bash
    cp .env.example .env
    ```

3.  **Install Dependencies:**
    This command will download the packages listed in your `go.mod` file.
    ```bash
    go mod tidy
    ```

### Running the Server

To start the server, run the following command from the `/backend` directory:
```bash
go run ./cmd/server/main.go