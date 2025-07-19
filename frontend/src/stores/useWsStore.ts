import apiService from "@/api"
import { create } from "zustand"
import { useGameStore } from "./useGameStore"

type WsStore = {
    ws: WebSocket | null
    status: "connecting" | "connected" | "disconnected" | "error"
    connect: () => void
    disconnect: () => void
}

export const useWsStore = create<WsStore>((set, get) => ({
    ws: null,
    status: "disconnected",
    connect: async () => {
        const { status, ws: existingWs } = get();
        
        // Don't start new connection if already connecting or connected
        if (status === "connecting" || status === "connected") {
            return;
        }

        // Close existing connection if any
        if (existingWs) {
            existingWs.close();
        }

        set({ status: "connecting" });

        try {
            const ws = await apiService.websocketConnect();
            set({ ws });

            ws.addEventListener("open", () => {
                set({ status: "connected" });
            });

            ws.addEventListener("close", () => {
                set({ ws: null, status: "disconnected" });
            });

            ws.addEventListener("error", () => {
                set({ ws: null, status: "error" });
            });

            ws.addEventListener("message", (event) => {
                try {
                    const msg = JSON.parse(event.data);
                    useGameStore.getState().handleWSMessage(msg);
                } catch {
                    // Ignore parsing errors
                }
            });
        } catch {
            set({ status: "error" });
        }
    },
    disconnect: () => {
        const { ws } = get();
        if (ws) {
            ws.close(1000, "User disconnect");
        }
        set({ ws: null, status: "disconnected" });
    },
}))
