import { UserWithToken } from "@/types";

export const login = async (username: string): Promise<UserWithToken> => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ username }),
    })
    return response.json()
}

export const refreshToken = async (refreshToken: string): Promise<UserWithToken> => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/refresh`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ refreshToken }),
    });
    
    if (!response.ok) {
        throw new Error('Failed to refresh token');
    }
    
    return response.json();
}

export const websocketConnect = (token: string) => {
    const wsUrl = process.env.NEXT_PUBLIC_API_URL?.replace('http', 'ws') || 'ws://localhost:8080';
    return new WebSocket(`${wsUrl}/ws?token=${encodeURIComponent(token)}`);
}
