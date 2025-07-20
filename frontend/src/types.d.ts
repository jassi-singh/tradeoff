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

interface PnlData {
    totalRealizedPnl: number;
    totalUnrealizedPnl: number;
}

interface PhaseData {
    phase: GamePhase;
    endTime: string;
}

interface CountData {
    longPositions: number;
    shortPositions: number;
    totalPlayers: number;
}

interface BasePlayerState {
    balance: number;
    activePosition: Position;
    closedPositions: ClosedPosition[];
}

interface Position {
    type: PositionType;
    entryPrice: number;
    entryTime: string; 
    pnl: number;
    pnlPercentage: number;
}

interface ClosedPosition extends Position {
    exitPrice: number;
    exitTime: string;
}

interface GameStateData extends PhaseData, CountData, PnlData, BasePlayerState {
    roundId: string;
    chartData: CandlestickData[];
}

interface PriceUpdateData {
    priceData: CandlestickData;
    updateLast: boolean;
}

export interface WebSocketMessage  {
    type:  "price_update" | "pnl_update" | "phase_update" | "count_update" | "game_state_sync"| "new_round";
    data: PriceUpdateData | PnLData | PhaseData | CountData | GameStateData 
}
