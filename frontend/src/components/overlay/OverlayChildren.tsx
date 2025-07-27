const OverlayChildren = () => {
  return (
    <div className="max-w-md w-full">
      {/* Header */}
      <div className="text-center mb-6">
        <h2 className="text-3xl font-bold text-white mb-2">TradeOff</h2>
        <p className="text-gray-400 text-sm">
          10-Minute Stock Market Simulator
        </p>
      </div>

      {/* Game Rules */}
      <div className="bg-gray-900/50 rounded-lg p-4 mb-6 border border-gray-700">
        <h3 className="text-lg font-semibold text-white mb-3">Game Rules</h3>
        <div className="space-y-3 text-sm">
          <div className="flex items-center gap-3">
            <div className="w-2 h-2 bg-blue-400 rounded-full"></div>
            <span className="text-gray-300">
              Start with{" "}
              <span className="text-green-400 font-semibold">$100</span> balance
              each round
            </span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-2 h-2 bg-green-400 rounded-full"></div>
            <span className="text-gray-300">
              Go <span className="text-green-400 font-semibold">Long</span> if
              you think price will rise
            </span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-2 h-2 bg-red-400 rounded-full"></div>
            <span className="text-gray-300">
              Go <span className="text-red-400 font-semibold">Short</span> if
              you think price will fall
            </span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-2 h-2 bg-yellow-400 rounded-full"></div>
            <span className="text-gray-300">
              One position per player at a time
            </span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-2 h-2 bg-purple-400 rounded-full"></div>
            <span className="text-gray-300">
              Real-time{" "}
              <span className="text-purple-400 font-semibold">Bitcoin/USD</span>{" "}
              price data
            </span>
          </div>
        </div>
      </div>

      {/* Join Form */}
      <div className="space-y-4">
        <div>
          <label
            htmlFor="username"
            className="block text-sm font-medium text-gray-300 mb-2"
          >
            Enter Your Trading Name
          </label>
          <input
            type="text"
            id="username"
            name="username"
            placeholder="e.g., CryptoTrader"
            className="w-full p-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 transition-all"
            required
          />
        </div>

        <button
          type="submit"
          className="w-full bg-green-600 text-white p-3 rounded-lg font-semibold hover:bg-green-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-green-500"
        >
          Join Trading Arena
        </button>
      </div>

      {/* Footer */}
      <div className="text-center mt-6">
        <p className="text-xs text-gray-500">
          Zero financial risk • Real-time multiplayer • Live leaderboard
        </p>
      </div>
    </div>
  );
};

export default OverlayChildren;
