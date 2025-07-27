"use client";

import { useGameStore } from "@/stores/useGameStore";
import useAuthStore from "@/stores/useAuthStore";

export default function Leaderboard() {
  const { leaderboardData, totalPlayers, longPositions, shortPositions } =
    useGameStore();
  const { user } = useAuthStore();

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const formatPercentage = (percentage: number) => {
    return `${percentage >= 0 ? "+" : ""}${percentage.toFixed(1)}%`;
  };

  // Show loading state if no leaderboard data yet
  if (!leaderboardData) {
    return (
      <div className="bg-gray-900 rounded-lg p-4 h-full">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-white">Leaderboard</h3>
          <div className="text-sm text-gray-400">
            {totalPlayers} players online
          </div>
        </div>
        <div className="flex justify-center items-center h-32">
          <div className="text-gray-400">Loading leaderboard...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-900 rounded-lg p-4 h-full">
      {/* Header */}
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold text-white">Leaderboard</h3>
        <div className="text-sm text-gray-400">
          {leaderboardData.totalPlayers} players online
        </div>
      </div>

      {/* Position Distribution */}
      <div className="flex justify-between text-sm mb-4">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 bg-green-500 rounded-full"></div>
          <span className="text-gray-400">
            Long: {leaderboardData.longPositions}
          </span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 bg-red-500 rounded-full"></div>
          <span className="text-gray-400">
            Short: {leaderboardData.shortPositions}
          </span>
        </div>
      </div>

      {/* Leaderboard List */}
      <div className="space-y-2 flex-1">
        {leaderboardData.players.slice(0, 8).map((player) => (
          <div
            key={player.playerId}
            className={`flex items-center justify-between p-2 rounded ${
              player.username === user?.username
                ? "bg-blue-900/30 border border-blue-500/50"
                : "bg-gray-800/50"
            }`}
          >
            <div className="flex items-center gap-2">
              <div className="flex items-center justify-center w-5 h-5 rounded-full bg-gray-700 text-xs font-bold">
                {player.rank}
              </div>
              <div className="flex flex-col">
                <span className="text-sm font-medium text-white">
                  {player.username}
                </span>
                {player.activePosition && (
                  <span
                    className={`text-xs ${
                      player.activePosition.type === "long"
                        ? "text-green-400"
                        : "text-red-400"
                    }`}
                  >
                    {player.activePosition.type.toUpperCase()}{" "}
                    {formatPercentage(player.activePosition.pnlPercentage)}
                  </span>
                )}
              </div>
            </div>
            <div className="text-right">
              <div className="text-sm font-mono text-white">
                {formatCurrency(player.balance)}
              </div>
              <div
                className={`text-xs ${
                  player.totalPnl >= 0 ? "text-green-400" : "text-red-400"
                }`}
              >
                {formatCurrency(player.totalPnl)}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Live Indicator */}
      <div className="flex items-center justify-center gap-2 text-xs text-green-400 mt-4">
        <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
        <span>Live Updates</span>
      </div>
    </div>
  );
}
