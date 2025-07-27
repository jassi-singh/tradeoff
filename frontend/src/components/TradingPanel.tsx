"use client";
import { useGameStore } from "@/stores/useGameStore";
import { PositionType } from "@/types";
import apiService from "@/api";

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

  return (
    <div className="bg-gray-900 rounded-lg p-6 space-y-6">
      {/* Balance Section */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-white">Balance</h3>
        <div className="space-y-2">
          <div className="flex justify-between">
            <span className="text-gray-400">Total Balance:</span>
            <span className="text-white font-mono">
              {formatCurrency(balance)}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-400">Realized P&L:</span>
            <span
              className={`font-mono ${
                totalRealizedPnl >= 0 ? "text-green-400" : "text-red-400"
              }`}
            >
              {formatCurrency(totalRealizedPnl)}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-400">Unrealized P&L:</span>
            <span
              className={`font-mono ${
                totalUnrealizedPnl >= 0 ? "text-green-400" : "text-red-400"
              }`}
            >
              {formatCurrency(totalUnrealizedPnl)}
            </span>
          </div>
        </div>
      </div>

      {/* Trading Buttons */}
      {phase === "live" && !activePosition && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-white">Trading</h3>
          <div className="space-y-3">
            <button
              onClick={() => handleTrade("long")}
              className="w-full bg-green-600 hover:bg-green-700 text-white font-semibold py-3 px-4 rounded-lg transition-colors"
            >
              Go Long
            </button>
            <button
              onClick={() => handleTrade("short")}
              className="w-full bg-red-600 hover:bg-red-700 text-white font-semibold py-3 px-4 rounded-lg transition-colors"
            >
              Go Short
            </button>
          </div>
        </div>
      )}

      {/* Close Position Button */}
      {phase === "live" && activePosition && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-white">Trading</h3>
          <div className="space-y-3">
            <button
              onClick={handleClosePosition}
              className="w-full bg-yellow-600 hover:bg-yellow-700 text-white font-semibold py-3 px-4 rounded-lg transition-colors"
            >
              Close Position
            </button>
          </div>
        </div>
      )}

      {/* Active Position */}
      {activePosition && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-white">Active Position</h3>
          <div className="bg-gray-800 rounded-lg p-4 space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-gray-400">Type:</span>
              <span
                className={`font-semibold ${
                  activePosition.type === "long"
                    ? "text-green-400"
                    : "text-red-400"
                }`}
              >
                {activePosition.type.toUpperCase()}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Entry Price:</span>
              <span className="text-white font-mono">
                {formatCurrency(activePosition.entryPrice)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Entry Time:</span>
              <span className="text-white text-sm">
                {new Date(activePosition.entryTime).toLocaleTimeString()}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">P&L:</span>
              <span
                className={`font-mono ${
                  activePosition.pnl >= 0 ? "text-green-400" : "text-red-400"
                }`}
              >
                {formatCurrency(activePosition.pnl)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">P&L %:</span>
              <span
                className={`font-mono ${
                  activePosition.pnlPercentage >= 0
                    ? "text-green-400"
                    : "text-red-400"
                }`}
              >
                {formatPercentage(activePosition.pnlPercentage)}
              </span>
            </div>
          </div>
        </div>
      )}

      {/* Game Phase Info */}
      <div className="space-y-2">
        <h3 className="text-lg font-semibold text-white">Game Status</h3>
        <div className="flex justify-between">
          <span className="text-gray-400">Phase:</span>
          <span
            className={`font-semibold ${
              phase === "live"
                ? "text-green-400"
                : phase === "closed"
                ? "text-red-400"
                : "text-yellow-400"
            }`}
          >
            {phase.toUpperCase()}
          </span>
        </div>
      </div>
    </div>
  );
}
