"use client";

import {
  createChart,
  IChartApi,
  CandlestickData,
  AreaData,
  ColorType,
  CandlestickSeries,
  AreaSeries,
} from "lightweight-charts";
import React, { useEffect, useRef } from "react";

const areaData = [
  { time: "2023-01-01", value: 100 },
  { time: "2023-01-02", value: 105 },
  { time: "2023-01-03", value: 102 },
];

const candlestickData = [
  { time: "2023-01-01", open: 100, high: 105, low: 95, close: 102 },
  { time: "2023-01-02", open: 105, high: 110, low: 102, close: 108 },
  { time: "2023-01-03", open: 102, high: 108, low: 98, close: 105 },
];

const ChartComponent: React.FC = () => {
  // Use a ref to attach to the container div
  const chartContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Ensure the ref is attached to a DOM element
    if (!chartContainerRef.current) {
      return;
    }

    const chartOptions = {
      layout: {
        textColor: "white",
        background: { type: ColorType.Solid, color: "black" },
      },
      width: chartContainerRef.current.clientWidth,
      height: 500, // Or make this a prop
    };

    // Create the chart instance
    const chart: IChartApi = createChart(
      chartContainerRef.current,
      chartOptions
    );

    // Add the candlestick series with its data
    const candlestickSeries = chart.addSeries(CandlestickSeries, {
      upColor: "#26a69a",
      downColor: "#ef5350",
      borderVisible: false,
      wickUpColor: "#26a69a",
      wickDownColor: "#ef5350",
    });
    candlestickSeries.setData(candlestickData);

    // Add the area series with its data
    const areaSeries = chart.addSeries(AreaSeries, {
      lineColor: "#2962FF",
      topColor: "#2962FF",
      bottomColor: "rgba(41, 98, 255, 0.28)",
    });
    areaSeries.setData(areaData);

    // Fit the content to the chart
    chart.timeScale().fitContent();

    // Make the chart responsive to window resizing
    const handleResize = () => {
      chart.applyOptions({ width: chartContainerRef.current?.clientWidth });
    };

    window.addEventListener("resize", handleResize);

    // Cleanup function to remove the chart and event listener on component unmount
    return () => {
      window.removeEventListener("resize", handleResize);
      chart.remove();
    };
  }, [areaData, candlestickData]); // Re-run the effect if data changes

  return (
    <div
      ref={chartContainerRef}
      className="chart-container"
      style={{ position: "relative" }}
    />
  );
};

export default ChartComponent;
