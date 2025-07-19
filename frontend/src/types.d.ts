export interface User {
    id: string;
    username: string;
}

export interface UserWithToken {
    user: User;
    token: string;
    refreshToken: string;
}

// Position types
export type PositionType = "long" | "short";

export interface Position {
    type: PositionType;
    entryPrice: number;
    entryTime: string; // ISO date string
    exitPrice?: number;
    exitTime?: string; // ISO date string
    profit?: number;
    profitPercentage?: number;
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

interface PositionUpdateMessage {
    type: "position_update";
    data: Position;
}

interface PositionClosedMessage {
    type: "position_closed";
    data: Position;
}

export type WebSocketMessage = PriceUpdateMessage | ChartDataMessage | RoundStatusMessage | GameStateMessage | PositionUpdateMessage | PositionClosedMessage;
