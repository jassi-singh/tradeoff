import { CandlestickData } from "lightweight-charts";
import { create } from "zustand";
import { GamePhase, WebSocketMessage } from "@/types";

type GameStore = {
    chartPriceData: CandlestickData[]
    phase: GamePhase
    phaseEndTime?: Date
    setChartPriceData: (chartPriceData: CandlestickData[]) => void
    appendSinglePriceData: (newData: CandlestickData) => void

    handleWSMessage: (msg: WebSocketMessage) => void
}

export const useGameStore = create<GameStore>((set, get) => ({
    chartPriceData: [],
    phaseEndTime: undefined,
    phase: "lobby",
    setChartPriceData: (chartPriceData: CandlestickData[]) => set({ chartPriceData }),
    appendSinglePriceData: (newData: CandlestickData) => {
        const chartPriceData = get().chartPriceData;
        if (!chartPriceData.length) {
            set((state) => ({
                chartPriceData: [...state.chartPriceData, newData]
            }));
            return;
        }
        const lastData = chartPriceData[chartPriceData.length - 1];

        if (typeof lastData.time !== 'number' || typeof newData.time !== 'number') {
            return;
        }

        const lastDate = new Date(lastData.time * 1000);
        const newDate = new Date(newData.time * 1000);

        if (lastDate.getUTCDate() === newDate.getUTCDate()) {
            chartPriceData.pop();
            lastData.high = Math.max(lastData.high, newData.high);
            lastData.low = Math.min(lastData.low, newData.low);
            lastData.close = newData.close;

            chartPriceData.push(lastData);
            set({ chartPriceData: [...chartPriceData] });
        } else {
            set((state) => ({
                chartPriceData: [...state.chartPriceData, newData]
            }));
        }
    },
    handleWSMessage: (msg: WebSocketMessage) => {
        const currentPhase = get().phase;
        switch (msg.type) {
            case "chart_data":
                set({ chartPriceData: msg.data.chartData });
                break;
            case "price_update":
                if (currentPhase === "live") {
                    get().appendSinglePriceData(msg.data);
                }
                break;
            case "round_status":
                set({ phaseEndTime: new Date(msg.data.nextPhaseTime), phase: msg.data.phase });
                break;

            case "game_state":
                set({
                    chartPriceData: msg.data.chartData,
                    phase: msg.data.phase,
                    phaseEndTime: msg.data.phaseEndTime ? new Date(msg.data.phaseEndTime) : undefined
                });
                break;
            default:
                console.warn(`Unhandled message type: ${(msg as WebSocketMessage).type}`);
        }
    }
}))
