export const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    minimumFractionDigits: 2,
  }).format(amount);
};

export const formatPercentage = (percentage: number): string => {
  return `${percentage >= 0 ? "+" : ""}${percentage.toFixed(2)}%`;
};

export const formatDuration = (seconds: number): string => {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, "0")}`;
};

export const formatTime = (seconds: number | null): string => {
  if (seconds === null) return "00:00";
  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${String(minutes).padStart(2, "0")}:${String(secs).padStart(2, "0")}`;
};

export const getPnlColor = (pnl: number): string => {
  if (pnl > 0) return "text-green-400";
  if (pnl < 0) return "text-red-400";
  return "text-gray-400";
};

export const getPhaseColor = (phase: string): string => {
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

export const getStatusColor = (status: string): string => {
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

export const getRiskLevel = (pnlPercentage: number): {
  level: "LOW" | "MED" | "HIGH";
  color: string;
} => {
  const absPercentage = Math.abs(pnlPercentage);
  if (absPercentage > 10) {
    return { level: "HIGH", color: "text-red-400" };
  } else if (absPercentage > 5) {
    return { level: "MED", color: "text-yellow-400" };
  } else {
    return { level: "LOW", color: "text-green-400" };
  }
}; 