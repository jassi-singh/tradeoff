import JoinGameOverlay from "@/components/overlay/JoinGameOverlay";
import OverlayChildren from "@/components/overlay/OverlayChildren";
import ChartComponent from "@/components/CandlestickChart";
import { GameInfo } from "@/components/GameInfo";
import TradingPanel from "@/components/TradingPanel";
import Leaderboard from "@/components/Leaderboard";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-8 bg-gradient-to-br from-black to-gray-900 text-white">
      <div className="w-full max-w-7xl flex-grow flex flex-col">
        <header className="mb-6">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-green-400 to-blue-400 bg-clip-text text-transparent">
            Trade Off
          </h1>
          <p className="text-gray-400 text-sm mt-2">
            10-Minute Stock Market Simulator
          </p>
        </header>

        {/* Top section: Chart and Leaderboard */}
        <div className="flex-grow flex gap-6 mb-6">
          <div className="flex-[2] relative bg-gray-900/80 backdrop-blur-sm rounded-xl border border-gray-700/30 p-4 shadow-xl">
            <GameInfo />
            <ChartComponent />
          </div>
          <div className="flex-[1]">
            <Leaderboard />
          </div>
        </div>

        {/* Bottom section: Trading Panel */}
        <div className="w-full">
          <TradingPanel />
        </div>
      </div>
      <JoinGameOverlay>
        <OverlayChildren />
      </JoinGameOverlay>
    </main>
  );
}
