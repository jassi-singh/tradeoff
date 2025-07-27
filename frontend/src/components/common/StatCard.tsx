interface StatCardProps {
  label: string;
  value: string | number;
  valueColor?: string;
  className?: string;
}

export const StatCard: React.FC<StatCardProps> = ({
  label,
  value,
  valueColor = "text-white",
  className = "",
}) => {
  return (
    <div className={`text-center ${className}`}>
      <div className="text-xs text-gray-400 mb-1">{label}</div>
      <div className={`font-mono text-sm ${valueColor}`}>
        {typeof value === "number" ? value.toFixed(2) : value}
      </div>
    </div>
  );
};
