"use client";

import { useGameStore } from "@/stores/useGameStore";
import { useWsStore } from "@/stores/useWsStore";
import { Timer } from "./Timer";
import { getPhaseColor, getStatusColor } from "@/utils/formatters";

export const GameInfo = () => {
  const { phase } = useGameStore();
  const { status } = useWsStore();

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
    <div className="flex justify-between items-center z-10 w-full absolute top-0 left-0 bg-black/20 backdrop-blur-sm p-6 text-white">
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2 text-sm">
          <div
            className={`w-2 h-2 rounded-full animate-pulse ${getPhaseColor(
              phase
            )}`}
          ></div>
          <span className="text-sm">
            Game Phase:{" "}
            <span className={`font-bold ${getPhaseColor(phase)}`}>
              {phase.toUpperCase()}
            </span>
          </span>
        </div>
        <div className="flex items-center gap-2 text-sm">
          <div
            className={`w-2 h-2 rounded-full animate-pulse ${getStatusColor(
              status
            )}`}
          ></div>
          <span className="text-sm">
            Status:{" "}
            <span className={`font-bold ${getStatusColor(status)}`}>
              {getStatusText(status)}
            </span>
          </span>
        </div>
      </div>
      <Timer />
    </div>
  );
};
