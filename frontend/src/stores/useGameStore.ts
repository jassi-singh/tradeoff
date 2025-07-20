import { CandlestickData } from "lightweight-charts";
import { create } from "zustand";
import { GamePhase, Position, ClosedPosition, WebSocketMessage, PnlData, PhaseData, CountData, GameStateData, PriceUpdateData } from "@/types";

type GameStore = {
    roundId: string
    chartData: CandlestickData[]
    phase:           GamePhase           
    endTime:         Date | null            
    longPositions:   number                     
    shortPositions:  number                     
    totalPlayers:    number                     
    balance:         number                
    activePosition:  Position | null
    closedPositions: ClosedPosition[] 
    totalRealizedPnl: number
    totalUnrealizedPnl: number

    // actions

    handleWSMessage: (msg: WebSocketMessage) => void
}

export const useGameStore = create<GameStore>((set, get) => ({
    roundId: "",
    chartData: [],
    endTime: null,
    phase: "lobby",
    longPositions: 0,
    shortPositions: 0,
    totalPlayers: 0,
    balance: 0,
    activePosition: null,
    closedPositions: [],
    totalRealizedPnl: 0,
    totalUnrealizedPnl: 0,
    
    
    setPnl: (pnlData: PnlData) => set({ totalRealizedPnl: pnlData.totalRealizedPnl, totalUnrealizedPnl: pnlData.totalUnrealizedPnl }),
    setPhase: (phaseData: PhaseData) => set({ phase: phaseData.phase, endTime: new Date(phaseData.endTime) }),
    setCount: (countData: CountData) => set({ longPositions: countData.longPositions, shortPositions: countData.shortPositions, totalPlayers: countData.totalPlayers }),
    setGameState: (gameStateData: GameStateData) => set({ roundId: gameStateData.roundId, chartData: gameStateData.chartData, phase: gameStateData.phase, endTime: new Date(gameStateData.endTime) }),
    handleWSMessage: (msg: WebSocketMessage) => {
        switch (msg.type) {
            case "price_update":
                const priceUpdate = msg.data as PriceUpdateData;
                set((state) => {
                    const newChartData = [...state.chartData];
                    if (priceUpdate.updateLast && newChartData.length > 0) {
                        // Update the last candle
                        newChartData[newChartData.length - 1] = priceUpdate.priceData;
                    } else {
                        // Append a new candle
                        newChartData.push(priceUpdate.priceData);
                    }
                    return { chartData: newChartData };
                });
                break;
            case "pnl_update": {
                const data = msg.data as PnlData;
                set({
                    totalRealizedPnl: data.totalRealizedPnl,
                    totalUnrealizedPnl: data.totalUnrealizedPnl,
                });
                break;
            }
            case "phase_update": {
                const data = msg.data as PhaseData;
                set({
                    phase: data.phase,
                    endTime: new Date(data.endTime),
                });
                break;
            }
            case "count_update": {
                const data = msg.data as CountData;
                set({
                    longPositions: data.longPositions,
                    shortPositions: data.shortPositions,
                    totalPlayers: data.totalPlayers,
                });
                break;
            }
            case "game_state_sync": 
            case "new_round": {
                const data = msg.data as GameStateData;
                set({
                    roundId: data.roundId,
                    chartData: data.chartData,
                    phase: data.phase,
                    endTime: new Date(data.endTime),
                    balance: data.balance,
                    activePosition: data.activePosition,
                    closedPositions: data.closedPositions,
                    totalRealizedPnl: data.totalRealizedPnl,
                    totalUnrealizedPnl: data.totalUnrealizedPnl,
                    longPositions: data.longPositions,
                    shortPositions: data.shortPositions,
                    totalPlayers: data.totalPlayers,
                });
                break;
            }
            default:
                console.warn(`Unhandled message type: ${(msg as WebSocketMessage).type}`);
        }
    }
}))
