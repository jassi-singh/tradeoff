"use client";
import { useGameStore } from "@/stores/useGameStore";
import { PositionType } from "@/types";
import apiService from "@/api";
import { useState, useEffect } from "react";

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

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const formatPercentage = (percentage: number) => {
    return `${percentage >= 0 ? "+" : ""}${percentage.toFixed(2)}%`;
  };

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, "0")}`;
  };

  const getPnlColor = (pnl: number) => {
    if (pnl > 0) return "text-green-400";
    if (pnl < 0) return "text-red-400";
    return "text-gray-400";
  };

  const getPhaseColor = (phase: string) => {
    switch (phase) {
      case "live":
        return "text-green-400";
      case "closed":
        return "text-red-400";
      case "lobby":
        return "text-yellow-400";
      default:
        return "text-gray-400";
    }
  };

  return (
    <div className="bg-gray-900 rounded-lg p-4">
      <div className="flex items-center justify-between gap-6">
        {/* Balance Section */}
        <div className="flex items-center gap-6">
          <div className="text-center">
            <div className="text-xs text-gray-400 mb-1">Balance</div>
            <div className="text-white font-mono text-lg font-bold">
              {formatCurrency(balance)}
            </div>
          </div>

          <div className="text-center">
            <div className="text-xs text-gray-400 mb-1">Realized P&L</div>
            <div
              className={`font-mono text-sm ${getPnlColor(totalRealizedPnl)}`}
            >
              {formatCurrency(totalRealizedPnl)}
            </div>
          </div>

          <div className="text-center">
            <div className="text-xs text-gray-400 mb-1">Unrealized P&L</div>
            <div
              className={`font-mono text-sm ${getPnlColor(totalUnrealizedPnl)}`}
            >
              {formatCurrency(totalUnrealizedPnl)}
            </div>
          </div>
        </div>

        {/* Trading Controls */}
        <div className="flex items-center gap-4">
          {phase === "live" && !activePosition && (
            <>
              <button
                onClick={() => handleTrade("long")}
                className="bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors text-sm"
              >
                LONG
              </button>
              <button
                onClick={() => handleTrade("short")}
                className="bg-red-600 hover:bg-red-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors text-sm"
              >
                SHORT
              </button>
            </>
          )}

          {phase === "live" && activePosition && (
            <button
              onClick={handleClosePosition}
              className="bg-yellow-600 hover:bg-yellow-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors text-sm"
            >
              CLOSE
            </button>
          )}
        </div>

        {/* Game Phase */}
        <div className="text-center">
          <div className="text-xs text-gray-400 mb-1">Phase</div>
          <div className={`font-semibold text-sm ${getPhaseColor(phase)}`}>
            {phase.toUpperCase()}
          </div>
        </div>

        {/* Active Position */}
        {activePosition && (
          <div className="flex items-center gap-6">
            <div className="text-center">
              <div className="text-xs text-gray-400 mb-1">Position</div>
              <div
                className={`text-sm font-semibold ${
                  activePosition.type === "long"
                    ? "text-green-400"
                    : "text-red-400"
                }`}
              >
                {activePosition.type.toUpperCase()}
              </div>
            </div>

            <div className="text-center">
              <div className="text-xs text-gray-400 mb-1">Entry Price</div>
              <div className="text-white font-mono text-sm">
                {formatCurrency(activePosition.entryPrice)}
              </div>
            </div>

            <div className="text-center">
              <div className="text-xs text-gray-400 mb-1">P&L</div>
              <div
                className={`font-mono text-sm ${getPnlColor(
                  activePosition.pnl
                )}`}
              >
                {formatCurrency(activePosition.pnl)}
              </div>
            </div>

            <div className="text-center">
              <div className="text-xs text-gray-400 mb-1">P&L %</div>
              <div
                className={`font-mono text-sm ${getPnlColor(
                  activePosition.pnlPercentage
                )}`}
              >
                {formatPercentage(activePosition.pnlPercentage)}
              </div>
            </div>

            <div className="text-center">
              <div className="text-xs text-gray-400 mb-1">Duration</div>
              <div className="text-white font-mono text-sm">
                {formatDuration(duration)}
              </div>
            </div>

            <div className="text-center">
              <div className="text-xs text-gray-400 mb-1">Risk</div>
              <div
                className={`text-xs font-semibold ${
                  Math.abs(activePosition.pnlPercentage) > 10
                    ? "text-red-400"
                    : Math.abs(activePosition.pnlPercentage) > 5
                    ? "text-yellow-400"
                    : "text-green-400"
                }`}
              >
                {Math.abs(activePosition.pnlPercentage) > 10
                  ? "HIGH"
                  : Math.abs(activePosition.pnlPercentage) > 5
                  ? "MED"
                  : "LOW"}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
