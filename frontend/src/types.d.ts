export interface User {
    id: string;
    username: string;
}

// New types for WebSocket messages
import { CandlestickData } from "lightweight-charts";

export type GamePhase = "lobby" | "live" | "closed";

interface PriceUpdateMessage {
    type: "price_update";
    data: CandlestickData;
}

interface ChartDataMessage {
    type: "chart_data";
    data: {
        chartData: CandlestickData[];
    };
}

interface RoundStatusMessage {
    type: "round_status";
    data: {
        phase: GamePhase;
        nextPhaseTime: string; // ISO date string
    };
}

interface GameStateMessage {
    type: "game_state";
    data: {
        phase: GamePhase;
        chartData: CandlestickData[];
        phaseEndTime?: string; // ISO date string
    };
}

export type WebSocketMessage = PriceUpdateMessage | ChartDataMessage | RoundStatusMessage | GameStateMessage;