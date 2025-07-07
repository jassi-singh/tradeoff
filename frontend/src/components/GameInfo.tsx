'use client';
import { useGameStore } from "@/stores/useGameStore";
import { Timer } from "./Timer";

export const GameInfo = () => {
  const phase = useGameStore((state) => state.phase);

  return (
    <div className="flex justify-between items-center z-10 w-full absolute top-0 left-0 bg-black/10 backdrop-blur-sm p-4 text-white">
      <PhaseDisplay phase={phase} />
      <Timer />
    </div>
  );
}

const PhaseDisplay: React.FC<{ phase: string }> = ({ phase }) => {
  const phaseText = {
    lobby: "Lobby",
    live: "Live",
    closed: "Closed",
  }[phase] || "Unknown Phase";

  const phaseStyle = {
    lobby: "bg-yellow-500",
    live: "bg-green-500 animate-ping",
    closed: "bg-red-500",
  }[phase] || "bg-gray-500";

  return (
    <div className="flex items-center">
      <span className={`inline-block mr-2 h-1.5 w-1.5 rounded-full ${phaseStyle}`} />
      {phaseText}
    </div>
  );
};
