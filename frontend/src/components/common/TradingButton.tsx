interface TradingButtonProps {
  type: "long" | "short" | "close";
  onClick: () => void;
  disabled?: boolean;
  className?: string;
}

export const TradingButton: React.FC<TradingButtonProps> = ({
  type,
  onClick,
  disabled = false,
  className = "",
}) => {
  const getButtonStyles = () => {
    switch (type) {
      case "long":
        return "bg-green-600 hover:bg-green-700 text-white";
      case "short":
        return "bg-red-600 hover:bg-red-700 text-white";
      case "close":
        return "bg-yellow-600 hover:bg-yellow-700 text-white";
      default:
        return "bg-gray-600 hover:bg-gray-700 text-white";
    }
  };

  const getButtonText = () => {
    switch (type) {
      case "long":
        return "LONG";
      case "short":
        return "SHORT";
      case "close":
        return "CLOSE";
      default:
        return "";
    }
  };

  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`font-semibold py-2 px-6 rounded-lg transition-colors text-sm ${getButtonStyles()} ${className} ${
        disabled ? "opacity-50 cursor-not-allowed" : ""
      }`}
    >
      {getButtonText()}
    </button>
  );
};
