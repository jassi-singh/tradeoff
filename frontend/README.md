# TradeOff - Frontend Service

![Next.js Badge](https://img.shields.io/badge/Next.js-000000?style=for-the-badge&logo=next.js&logoColor=white)
![React Badge](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![Tailwind CSS Badge](https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)

This is the frontend client for **TradeOff**, a real-time, massively multiplayer stock market simulator.

This application is built with Next.js and React. It is responsible for rendering the game interface, handling user interactions, and communicating with the Go backend service via REST APIs and WebSockets.

---

## Project Status: Phase 1 (Initial Scaffolding)

The project is in its initial setup phase. The Next.js application has been scaffolded using the App Router, and the basic folder structure is in place.

The immediate goal is to build the core UI components and establish the initial connection with the backend for the guest identity system.

---

## Technology Stack

* **Framework:** [Next.js](https://nextjs.org/) (with App Router)
* **Library:** [React](https://react.dev/)
* **Styling:** [Tailwind CSS](https://tailwindcss.com/)
* **State Management (Planned):** Zustand or React Context
* **Charting Library (Planned):** TradingView Lightweight Charts

---

## Getting Started

Follow these instructions to get the frontend development server running on your local machine.

### Prerequisites

* [Node.js](https://nodejs.org/en) (version 18.17 or newer)
* [pnpm](https://pnpm.io/installation) (recommended package manager)

### Installation & Setup

1.  **Clone the Repository:**
    This frontend lives in the `/frontend` directory of the main monorepo.
    ```bash
    git clone [https://github.com/jassi-singh/tradeoff.git](https://github.com/jassi-singh/tradeoff.git)
    cd tradeoff/frontend
    ```

2.  **Set Up Environment Variables:**
    Copy the example environment file. This will be used to store the URL of your backend server.
    ```bash
    cp .env.example .env
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