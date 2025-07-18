import JoinGameOverlay from "@/components/overlay/JoinGameOverlay";
import OverlayChildren from "@/components/overlay/OverlayChildren";
import ChartComponent from "@/components/CandlestickChart";
import { GameInfo } from "@/components/GameInfo";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-8 bg-black text-white">
      <div className="w-full max-w-7xl flex-grow flex flex-col">
        <header className="mb-8">
          <h1 className="text-3xl font-bold">TradeOff</h1>
          <p className="text-gray-400">Minimalist Trading Dashboard</p>
        </header>
        <div className="flex-grow flex gap-4">
          <div className="flex-[2] relative">
            <GameInfo />
            <ChartComponent />
          </div>
          <div className="flex-[1] flex flex-col gap-4">
          </div>
        </div>
      </div>
      <JoinGameOverlay>
        <OverlayChildren />
      </JoinGameOverlay>
    </main>
  );
}
