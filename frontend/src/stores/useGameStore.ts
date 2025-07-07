import { CandlestickData } from "lightweight-charts";
import { create } from "zustand";

type GameStore = {
    chartPriceData: CandlestickData[]
    phase: "live" | "lobby" | "closed"
    phaseEndTime?: Date
    setChartPriceData: (chartPriceData: CandlestickData[]) => void
    appendSinglePriceData: (newData: CandlestickData) => void

    handleWSMessage: (msg: { type: string, data: any }) => void
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

        // if both has same day update last else add new
        const lastDate = new Date((lastData.time as any) * 1000);
        const newDate = new Date((newData.time as any) * 1000);

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
    handleWSMessage: (msg: { type: string, data: any }) => {
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
                console.warn(`Unhandled message type: ${msg.type}`);
        }
    }
}))
