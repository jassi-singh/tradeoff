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

## Current Project Status (As of June 29, 2025)

The project is currently in the initial stages of Phase 1. The foundational structure for the Go backend and Next.js frontend has been established within a monorepo, and the local development environment is containerized with Docker.

---

## The Vision: A Three-Phase Architectural Journey

This project is planned to evolve through three distinct phases, demonstrating a professional approach to building and scaling a modern web application.

### Phase 1: The Single-Player Core
* **Goal:** Build a functional, engaging, single-player version of the game. The core loop, trading mechanics, and UI must be solid.
* **Architecture:** Monolith (Single Go Backend + Next.js Frontend).
* **Status:** 🚧 **In Progress**

### Phase 2: The Multiplayer MVP
* **Goal:** Evolve the game into a complete, shareable, multiplayer experience with live public leaderboards, scheduled rounds, and viral features like the P&L Share Card.
* **Architecture:** The application will still be a **Monolith**, enhanced to handle concurrent players and social features.
* **Status:** ⏳ **Planned**

### Phase 3: The Scalability Case Study
* **Goal:** Use professional load-testing tools (`k6`) to prove the monolith's performance limits under simulated multiplayer load, and then evolve the architecture to solve those bottlenecks.
* **Architecture:** Refactor the system into an **Evolved Microservices Architecture**. Key components (e.g., data ingestion, trade processing) will be extracted into separate, scalable services.
* **Status:** ⏳ **Planned**

*[This README will be updated with "Before vs. After" performance benchmark graphs upon completion of Phase 3.]*

---

## Tech Stack Overview

| Component     | Technology                                           |
| :------------ | :--------------------------------------------------- |
| **Frontend** | Next.js, React, TypeScript, Tailwind CSS             |
| **Backend** | Go, Chi (Router), Gorilla WebSocket                  |
| **Database** | PostgreSQL                                         |
| **Orchestration** | Docker, Docker Compose                           |
| **Planned** | Redis (Phase 2), NATS/RabbitMQ (Phase 3)             |

---

## Getting Started

This project is containerized with Docker for a consistent and easy-to-manage development environment.

### Prerequisites

* [Git](https://git-scm.com/downloads)
* [Docker](https://www.docker.com/products/docker-desktop/) & Docker Compose

### Local Development Setup

1.  **Clone the Repository:**
    ```bash
    git clone [https://github.com/jassi-singh/tradeoff.git](https://github.com/jassi-singh/tradeoff.git)
    cd tradeoff
    ```

2.  **Configure Environment:**
    * Copy the example environment file for the backend: `cp backend/.env.example backend/.env`
    * Copy the example environment file for the frontend: `cp frontend/.env.local.example frontend/.env.local`
    * *(No edits are needed yet, but these files must exist for Docker Compose to work)*

3.  **Run the Application:**
    Use Docker Compose to build the images and start the services.
    ```bash
    docker-compose up --build
    ```
    * The frontend will be available at `http://localhost:3000`.
    * The backend will be available at `http://localhost:8080`.

### Alternative: Running Services Manually

If you prefer not to use Docker, you can run each service manually. For detailed instructions, see the README file within its directory:
* **[Backend README](./backend/README.md)**
* **[Frontend README](./frontend/README.md)**