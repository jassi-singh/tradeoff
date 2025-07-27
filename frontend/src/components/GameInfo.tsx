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

  const getStatusColorType = (
    status: string
  ): "success" | "warning" | "danger" | "default" => {
    switch (status) {
      case "connected":
        return "success";
      case "connecting":
        return "warning";
      case "error":
        return "danger";
      case "disconnected":
        return "default";
      default:
        return "default";
    }
  };

  const getPhaseColorType = (
    phase: string
  ): "success" | "warning" | "danger" | "default" => {
    switch (phase) {
      case "live":
        return "success";
      case "lobby":
        return "warning";
      case "closed":
        return "danger";
      default:
        return "default";
    }
  };

  const getColorClass = (type: string) => {
    switch (type) {
      case "success":
        return "text-green-400";
      case "warning":
        return "text-yellow-400";
      case "danger":
        return "text-red-400";
      default:
        return "text-gray-400";
    }
  };

  return (
    <div className="flex justify-between items-center z-10 w-full absolute top-0 left-0 bg-black/20 backdrop-blur-sm p-6 text-white">
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2 text-sm">
          <div
            className={`w-2 h-2 rounded-full animate-pulse ${
              getPhaseColorType(phase) === "success"
                ? "bg-green-400"
                : getPhaseColorType(phase) === "warning"
                ? "bg-yellow-400"
                : getPhaseColorType(phase) === "danger"
                ? "bg-red-400"
                : "bg-gray-400"
            }`}
          ></div>
          <span className="text-sm">
            Game Phase:{" "}
            <span
              className={`font-bold ${getColorClass(getPhaseColorType(phase))}`}
            >
              {phase.toUpperCase()}
            </span>
          </span>
        </div>
        <div className="flex items-center gap-2 text-sm">
          <div
            className={`w-2 h-2 rounded-full animate-pulse ${
              getStatusColorType(status) === "success"
                ? "bg-green-400"
                : getStatusColorType(status) === "warning"
                ? "bg-yellow-400"
                : getStatusColorType(status) === "danger"
                ? "bg-red-400"
                : "bg-gray-400"
            }`}
          ></div>
          <span className="text-sm">
            Status:{" "}
            <span
              className={`font-bold ${getColorClass(
                getStatusColorType(status)
              )}`}
            >
              {getStatusText(status)}
            </span>
          </span>
        </div>
      </div>
      <Timer />
    </div>
  );
};
