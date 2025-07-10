import { UserWithToken } from "@/types";

export const login = async (username: string): Promise<UserWithToken> => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/login`, {
        method: "POST",
        body: JSON.stringify({ username }),
    })
    return response.json()
}

export const websocketConnect = (playerId: string) => {
    return new WebSocket(`${process.env.NEXT_PUBLIC_API_URL}/ws?playerId=${playerId}`)
}
