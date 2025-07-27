export interface User {
    id: string;
    username: string;
}

export interface UserWithToken {
    user: User;
    token: string;
    refreshToken: string;
}

// New types for WebSocket messages
import { CandlestickData } from "lightweight-charts";

export type GamePhase = "lobby" | "live" | "closed";
export type PositionType = "long" | "short";

export interface PnlData {
    pnl: number;
    activePnl: number;
    activePnlPercentage: number;
    balance: number;
}

export interface PhaseData {
    phase: GamePhase;
    endTime: string;
}

export interface CountData {
    longPositions: number;
    shortPositions: number;
    totalPlayers: number;
}

export interface BasePlayerState {
    balance: number;
    activePosition: Position | null;
    closedPositions: ClosedPosition[];
}

export interface Position {
    type: PositionType;
    entryPrice: number;
    entryTime: string; 
    pnl: number;
    pnlPercentage: number;
}

export interface ClosedPosition extends Position {
    exitPrice: number;
    exitTime: string;
}

export interface GameStateData extends PhaseData, CountData, PnlData, BasePlayerState {
    roundId: string;
    chartData: CandlestickData[];
}

export interface PriceUpdateData {
    priceData: CandlestickData;
    updateLast: boolean;
}

export interface WebSocketMessage  {
    type:  "price_update" | "pnl_update" | "phase_update" | "count_update" | "game_state_sync"| "new_round";
    data: PriceUpdateData | PnlData | PhaseData | CountData | GameStateData 
}
