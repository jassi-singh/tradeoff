interface LoadingSpinnerProps {
  message?: string;
  className?: string;
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  message = "Loading...",
  className = "",
}) => {
  return (
    <div className={`flex justify-center items-center ${className}`}>
      <div className="text-gray-400">{message}</div>
    </div>
  );
};
