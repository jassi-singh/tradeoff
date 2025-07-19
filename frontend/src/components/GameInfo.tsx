"use client";

import { useGameStore } from "@/stores/useGameStore";
import { useWsStore } from "@/stores/useWsStore";
import { Timer } from "./Timer";

export const GameInfo = () => {
  const { phase } = useGameStore();
  const { status } = useWsStore();

  const getPhaseColor = (phase: string) => {
    switch (phase) {
      case "lobby":
        return "text-yellow-400";
      case "live":
        return "text-green-400";
      case "closed":
        return "text-red-400";
      default:
        return "text-gray-400";
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "connected":
        return "text-green-400";
      case "connecting":
        return "text-yellow-400";
      case "error":
        return "text-red-400";
      case "disconnected":
        return "text-gray-400";
      default:
        return "text-gray-400";
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case "connected":
        return "Connected";
      case "connecting":
        return "Connecting...";
      case "error":
        return "Connection Error";
      case "disconnected":
        return "Disconnected";
      default:
        return "Unknown";
    }
  };

  return (
    <div className="mb-4 p-4 bg-gray-800 rounded-lg">
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-4">
          <div className="text-sm">
            Game Phase:{" "}
            <span className={`font-bold ${getPhaseColor(phase)}`}>
              {phase.toUpperCase()}
            </span>
          </div>
          <div className="text-sm">
            Status:{" "}
            <span className={`font-bold ${getStatusColor(status)}`}>
              {getStatusText(status)}
            </span>
          </div>
        </div>
        <Timer />
      </div>
    </div>
  );
};
