import { CandlestickData } from "lightweight-charts";
import { create } from "zustand";

type GameStore = {
    chartPriceData: CandlestickData[]
    phaseEndTime?: Date
    setChartPriceData: (chartPriceData: CandlestickData[]) => void
    appendSinglePriceData: (newData: CandlestickData) => void

    handleWSMessage: (msg: { type: string, data: any }) => void
}

export const useGameStore = create<GameStore>((set, get) => ({
    chartPriceData: [],
    phaseEndTime: undefined,
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
        if (msg.type === "price_update") {
            get().appendSinglePriceData(msg.data);
        }

        if (msg.type === "round_status") {
            if (msg.data.phase === "live")
                get().setChartPriceData(msg.data.chartData);

            set({ phaseEndTime: new Date(msg.data.nextPhaseTime) });
        }
    }
}))
