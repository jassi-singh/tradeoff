import { websocketConnect } from "@/api"
import { create } from "zustand"
import { useGameStore } from "./useGameStore"

type WsStore = {
    ws: WebSocket | null
    status: "connecting" | "connected" | "disconnected" | "error"
    connect: (playerId: string) => void
    disconnect: () => void
}

export const useWsStore = create<WsStore>((set) => ({
    ws: null,
    status: "disconnected",
    connect: async (playerId: string) => {
        set({ status: "connecting" })
        const ws = websocketConnect(playerId)

        ws.addEventListener("open", () => {
            set({ status: "connected" })
        })
        ws.addEventListener("close", () => {
            set({ ws: null, status: "disconnected" })
        })
        ws.addEventListener("error", () => {
            set({ ws: null, status: "error" })
        })

        ws.addEventListener("message", (event) => {
            const msg = JSON.parse(event.data);
            useGameStore.getState().handleWSMessage(msg);
        });
    },
    disconnect: () => {
        set({ ws: null, status: "disconnected" })
    },
}))
