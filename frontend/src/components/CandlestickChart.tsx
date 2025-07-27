"use client";

import { useGameStore } from "@/stores/useGameStore";
import {
  createChart,
  IChartApi,
  ColorType,
  ISeriesApi,
  CandlestickSeries,
  DeepPartial,
  ChartOptions,
} from "lightweight-charts";
import React, { useEffect, useRef } from "react";

const ChartComponent: React.FC = () => {
  const { chartData: candlestickData } = useGameStore();
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);
  const seriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);

  useEffect(() => {
    if (!chartContainerRef.current) {
      return;
    }

    const chartOptions: DeepPartial<ChartOptions> = {
      layout: {
        textColor: "#ffffff",
        background: { type: ColorType.Solid, color: "#0a0a0a" },
      },
      grid: {
        vertLines: {
          color: "#374151",
        },
        horzLines: {
          color: "#374151",
        },
      },
      width: chartContainerRef.current.clientWidth,
    };

    const chart = createChart(chartContainerRef.current, chartOptions);
    chartRef.current = chart;

    const candlestickSeries = chart.addSeries(CandlestickSeries, {
      upColor: "#00d4aa",
      downColor: "#ff6b6b",
      borderVisible: false,
      wickUpColor: "#00d4aa",
      wickDownColor: "#ff6b6b",
    });
    seriesRef.current = candlestickSeries;

    // Add the area series with its data
    // const areaSeries = chart.addSeries(AreaSeries, {
    //   lineColor: "#2962FF",
    //   topColor: "#2962FF",
    //   bottomColor: "rgba(41, 98, 255, 0.28)",
    // });
    // areaSeries.setData(areaData);

    // Fit the content to the chart
    chart.timeScale().fitContent();

    const handleResize = () => {
      if (chartRef.current && chartContainerRef.current) {
        chartRef.current.applyOptions({
          width: chartContainerRef.current.clientWidth,
        });
      }
    };

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      if (chartRef.current) {
        chartRef.current.remove();
      }
    };
  }, []);

  useEffect(() => {
    if (seriesRef.current && candlestickData) {
      seriesRef.current.setData(candlestickData);
    }
  }, [candlestickData]);

  return (
    <div
      ref={chartContainerRef}
      className="chart-container h-full bg-gray-900/80 backdrop-blur-sm rounded-xl border border-gray-700/30"
      style={{ position: "relative" }}
    />
  );
};

export default ChartComponent;
