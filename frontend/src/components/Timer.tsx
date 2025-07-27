"use client";
import { useGameStore } from "@/stores/useGameStore";
import { useEffect, useState } from "react";
import { formatTime } from "@/utils/formatters";

export const Timer = () => {
  const endTime = useGameStore((state) => state.endTime);
  const [timeLeft, setTimeLeft] = useState<number | null>(null);

  useEffect(() => {
    if (!endTime) {
      setTimeLeft(null);
      return;
    }

    const calculateTimeLeft = () => {
      const now = new Date();
      const difference = endTime.getTime() - now.getTime();
      return Math.max(0, Math.floor(difference / 1000));
    };

    setTimeLeft(calculateTimeLeft());

    const interval = setInterval(() => {
      setTimeLeft(calculateTimeLeft());
    }, 1000);

    return () => clearInterval(interval);
  }, [endTime]);

  return (
    <div className="text-center">
      <div className="text-2xl font-bold text-white font-mono">
        {formatTime(timeLeft)}
      </div>
    </div>
  );
};
