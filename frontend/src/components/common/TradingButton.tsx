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
  const getButtonClass = () => {
    switch (type) {
      case "long":
        return "bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white";
      case "short":
        return "bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white";
      case "close":
        return "bg-gradient-to-r from-yellow-500 to-yellow-600 hover:from-yellow-600 hover:to-yellow-700 text-black";
      default:
        return "bg-gray-600 hover:bg-gray-700 text-white border border-gray-500";
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
      className={`font-semibold py-2 px-6 rounded-lg transition-all duration-200 text-sm ${getButtonClass()} ${className} ${
        disabled
          ? "opacity-50 cursor-not-allowed"
          : "hover:transform hover:scale-105"
      }`}
    >
      {getButtonText()}
    </button>
  );
};
