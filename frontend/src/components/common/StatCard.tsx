interface StatCardProps {
  label: string;
  value: string | number;
  valueColor?: "success" | "danger" | "warning" | "info" | "default";
  className?: string;
}

export const StatCard: React.FC<StatCardProps> = ({
  label,
  value,
  valueColor = "default",
  className = "",
}) => {
  const getValueColorClass = () => {
    switch (valueColor) {
      case "success":
        return "text-green-400";
      case "danger":
        return "text-red-400";
      case "warning":
        return "text-yellow-400";
      case "info":
        return "text-blue-400";
      default:
        return "text-white";
    }
  };

  return (
    <div className={`text-center ${className}`}>
      <div className="text-xs text-gray-400 mb-1 font-medium">{label}</div>
      <div
        className={`font-mono text-sm font-semibold ${getValueColorClass()}`}
      >
        {typeof value === "number" ? value.toFixed(2) : value}
      </div>
    </div>
  );
};
