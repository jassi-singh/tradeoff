"use client";
import { useGameStore } from "@/stores/useGameStore";
import { useState, useEffect } from "react";
import { StatCard } from "./common/StatCard";
import { TradingButton } from "./common/TradingButton";
import {
  formatCurrency,
  formatPercentage,
  formatDuration,
  getPnlColor,
  getPhaseColor,
  getRiskLevel,
} from "@/utils/formatters";

export default function TradingPanel() {
  const {
    balance,
    totalRealizedPnl,
    totalUnrealizedPnl,
    activePosition,
    phase,
    handleTrade,
    handleClosePosition,
  } = useGameStore();

  const [duration, setDuration] = useState(0);

  // Update duration every second when there's an active position
  useEffect(() => {
    if (!activePosition) {
      setDuration(0);
      return;
    }

    const interval = setInterval(() => {
      const entryTime = new Date(activePosition.entryTime).getTime();
      const now = Date.now();
      setDuration(Math.floor((now - entryTime) / 1000));
    }, 1000);

    return () => clearInterval(interval);
  }, [activePosition]);

  return (
    <div className="bg-gray-900/80 backdrop-blur-sm rounded-xl border border-gray-700/30 p-6 shadow-xl">
      <div className="flex items-center justify-between gap-6">
        {/* Balance Section */}
        <div className="flex items-center gap-6">
          <StatCard
            label="Balance"
            value={formatCurrency(balance)}
            valueColor="default"
            className="text-lg font-bold"
          />

          <StatCard
            label="Realized P&L"
            value={formatCurrency(totalRealizedPnl)}
            valueColor={getPnlColor(totalRealizedPnl)}
          />

          <StatCard
            label="Unrealized P&L"
            value={formatCurrency(totalUnrealizedPnl)}
            valueColor={getPnlColor(totalUnrealizedPnl)}
          />
        </div>

        {/* Trading Controls */}
        <div className="flex items-center gap-4">
          {phase === "live" && !activePosition && (
            <>
              <TradingButton type="long" onClick={() => handleTrade("long")} />
              <TradingButton
                type="short"
                onClick={() => handleTrade("short")}
              />
            </>
          )}

          {phase === "live" && activePosition && (
            <TradingButton type="close" onClick={handleClosePosition} />
          )}
        </div>

        {/* Game Phase */}
        <StatCard
          label="Phase"
          value={phase.toUpperCase()}
          valueColor={getPhaseColor(phase)}
          className="font-semibold"
        />

        {/* Active Position */}
        {activePosition && (
          <div className="flex items-center gap-6">
            <StatCard
              label="Position"
              value={activePosition.type.toUpperCase()}
              valueColor={activePosition.type === "long" ? "success" : "danger"}
              className="font-semibold"
            />

            <StatCard
              label="Entry Price"
              value={formatCurrency(activePosition.entryPrice)}
            />

            <StatCard
              label="P&L"
              value={formatCurrency(activePosition.pnl)}
              valueColor={getPnlColor(activePosition.pnl)}
            />

            <StatCard
              label="P&L %"
              value={formatPercentage(activePosition.pnlPercentage)}
              valueColor={getPnlColor(activePosition.pnlPercentage)}
            />

            <StatCard label="Duration" value={formatDuration(duration)} />

            <StatCard
              label="Risk"
              value={getRiskLevel(activePosition.pnlPercentage).level}
              valueColor={
                getRiskLevel(activePosition.pnlPercentage).level === "HIGH"
                  ? "danger"
                  : getRiskLevel(activePosition.pnlPercentage).level === "MED"
                  ? "warning"
                  : "success"
              }
              className="text-xs font-semibold"
            />
          </div>
        )}
      </div>
    </div>
  );
}
