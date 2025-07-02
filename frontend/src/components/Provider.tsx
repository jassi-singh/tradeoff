"use client";

import useAuthStore from "@/stores/useAuthStore";
import { useChartStore } from "@/stores/useChartStore";
import { useWsStore } from "@/stores/useWsStore";
import { CandlestickData } from "lightweight-charts";
import { useEffect } from "react";

const Provider = ({ children }: { children: React.ReactNode }) => {
  const { user } = useAuthStore();
  const {ws, connect, disconnect } = useWsStore();
  const { setChartPriceData } = useChartStore();

  useEffect(() => {
    if (user) {
      connect(user.id);
    }

    return () => {
      disconnect();
    };
  }, [user]);

  useEffect(() => {
    if (ws) {
      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === "price_data") {
          setChartPriceData(data.payload as CandlestickData[]);
        }
      }
    }
  }, [ws]);

  return <>{children}</>;
};

export default Provider;
