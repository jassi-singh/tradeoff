import JoinGameOverlay from "@/components/overlay/JoinGameOverlay";
import OverlayChildren from "@/components/overlay/OverlayChildren";
import ChartComponent from "@/components/CandlestickChart";
import { GameInfo } from "@/components/GameInfo";
import TradingPanel from "@/components/TradingPanel";
import Leaderboard from "@/components/Leaderboard";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-8 bg-black text-white">
      <div className="w-full max-w-7xl flex-grow flex flex-col">
        <header className="mb-4">
          <h1 className="text-3xl font-bold">TradeOff</h1>
        </header>

        {/* Top section: Chart and Leaderboard */}
        <div className="flex-grow flex gap-4 mb-4">
          <div className="flex-[2] relative">
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
