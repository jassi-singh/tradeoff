import { websocketConnect } from "@/api"
import { create } from "zustand"

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
        set({ ws, status: "connected" })
    },
    disconnect: () => {
        set({ ws: null, status: "disconnected" })
    },
}))