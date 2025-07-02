import { User } from "@/types";

export const createUser = async (username: string): Promise<User> => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/player`, {
        method: "POST",
        body: JSON.stringify({ username }),
    })
    return response.json()
}

export const websocketConnect = (playerId: string) => {
    return new WebSocket(`${process.env.NEXT_PUBLIC_API_URL}/ws?playerId=${playerId}`)
}