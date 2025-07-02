import { CandlestickData } from "lightweight-charts";
import { create } from "zustand";

type ChartStore = {
    chartPriceData: CandlestickData[]
    setChartPriceData: (chartPriceData: CandlestickData[]) => void
}

export const useChartStore = create<ChartStore>((set) => ({
    chartPriceData: [],
    setChartPriceData: (chartPriceData: CandlestickData[]) => set({ chartPriceData })
}))