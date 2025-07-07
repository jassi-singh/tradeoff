import { useWsStore } from "@/stores/useWsStore";
import { WebSocketMessage } from "@/types";
import { useEffect } from "react";

export const useWS = (callback: (msg: WebSocketMessage) => void) => {
    const ws = useWsStore((state) => state.ws);
    const status = useWsStore((state) => state.status);

    useEffect(() => {
        if (!ws || status !== "connected") {
            return
        };

        const messageHandler = (event: MessageEvent) => {
            const msg = JSON.parse(event.data) as WebSocketMessage;
            callback(msg);
        };

        ws.addEventListener("message", messageHandler);

        return () => {
            ws.removeEventListener("message", messageHandler);
        };
    }, [ws, status, callback]);

}
