import { websocketConnect } from "@/api"
import { create } from "zustand"
import { useGameStore } from "./useGameStore"
import useAuthStore from "./useAuthStore"

type WsStore = {
    ws: WebSocket | null
    status: "connecting" | "connected" | "disconnected" | "error"
    connect: (token?: string) => void
    disconnect: () => void
}

export const useWsStore = create<WsStore>((set, get) => ({
    ws: null,
    status: "disconnected",
    connect: async (providedToken?: string) => {
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

        let token: string | null = providedToken || null;
        
        // Get token from auth store if not provided
        if (!token) {
            const authState = useAuthStore.getState();
            token = authState.token;
            
            // Check if token is expired and refresh if needed
            if (!token || authState.isTokenExpired(token)) {
                token = await authState.refreshAuthToken();
                
                if (!token) {
                    set({ status: "error" });
                    return;
                }
            }
        }

        if (!token) {
            set({ status: "error" });
            return;
        }

        try {
            const ws = websocketConnect(token);
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
