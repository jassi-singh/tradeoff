"use client";

import { useGameStore } from "@/stores/useGameStore";
import useAuthStore from "@/stores/useAuthStore";
import { LoadingSpinner } from "./common/LoadingSpinner";
import { formatCurrency } from "@/utils/formatters";

export default function Leaderboard() {
  const { leaderboardData, totalPlayers, longPositions, shortPositions } =
    useGameStore();
  const { user } = useAuthStore();

  // Show loading state if no leaderboard data yet
  if (!leaderboardData) {
    return (
      <div className="bg-gray-900/80 backdrop-blur-sm rounded-xl border border-gray-700/30 p-6 h-full shadow-xl">
        <div className="flex justify-between items-center mb-6">
          <h3 className="text-xl font-semibold text-white">Leaderboard</h3>
          <div className="flex items-center gap-2 text-sm">
            <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
            <span className="text-gray-400">{totalPlayers} players online</span>
          </div>
        </div>
        <LoadingSpinner message="Loading leaderboard..." className="h-32" />
      </div>
    );
  }

  return (
    <div className="bg-gray-900/80 backdrop-blur-sm rounded-xl border border-gray-700/30 p-6 h-full shadow-xl">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-xl font-semibold text-white">Leaderboard</h3>
        <div className="flex items-center gap-2 text-sm">
          <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
          <span className="text-gray-400">{totalPlayers} players online</span>
        </div>
      </div>

      {/* Position Distribution */}
      <div className="flex justify-between text-sm mb-6">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 bg-green-500 rounded-full"></div>
          <span className="text-gray-400">Long: {longPositions}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 bg-red-500 rounded-full"></div>
          <span className="text-gray-400">Short: {shortPositions}</span>
        </div>
      </div>

      {/* Leaderboard List */}
      <div className="space-y-3 flex-1">
        {leaderboardData.map((player, index) => (
          <div
            key={player.playerId}
            className={`bg-gray-800/50 backdrop-blur-sm rounded-lg p-4 transition-all duration-200 border ${
              player.playerId === user?.id
                ? "border-blue-500/50 bg-blue-900/30"
                : "border-gray-700/30 hover:border-gray-500/50"
            }`}
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="flex items-center justify-center w-6 h-6 rounded-full bg-gray-700 text-xs font-bold text-white">
                  {index + 1}
                </div>
                <div className="flex flex-col">
                  <span className="text-sm font-medium text-white">
                    {player.username}
                  </span>
                </div>
              </div>
              <div className="text-right">
                <div className="text-sm font-mono text-white">
                  {formatCurrency(player.activeBalance)}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Live Indicator */}
      <div className="flex items-center justify-center gap-2 mt-6">
        <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
        <span className="text-green-400 text-xs">Live Updates</span>
      </div>
    </div>
  );
}
