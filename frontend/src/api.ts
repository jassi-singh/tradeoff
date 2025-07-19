import { UserWithToken, PositionType, Position } from "@/types";

class ApiService {
    private apiUrl: string;
    private token: string;

    constructor() {
        this.apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
        this.token = "";
    }

    setToken(token: string) {
        this.token = token;
    }

    async login(username: string): Promise<UserWithToken> {
        const response = await fetch(`${this.apiUrl}/api/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username }),
        });

        if (!response.ok) {
            throw new Error('Failed to login');
        }
        return response.json();
    }

    async refreshToken(refreshToken: string): Promise<UserWithToken> {
        const response = await fetch(`${this.apiUrl}/api/refresh`, {
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

    async websocketConnect(): Promise<WebSocket> {
        const wsUrl = this.apiUrl?.replace('http', 'ws') || 'ws://localhost:8080';
        return new WebSocket(`${wsUrl}/ws?token=${encodeURIComponent(this.token)}`);
    }

    async createPosition(positionType: PositionType): Promise<Position> {
        const response = await fetch(`${this.apiUrl}/api/position`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${this.token}`,
            },
            body: JSON.stringify({ type: positionType }),
        });

        if (!response.ok) {
            throw new Error('Failed to create position');
        }

        return response.json();
    }

    async closePosition(): Promise<void> {
        const response = await fetch(`${this.apiUrl}/api/close-position`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${this.token}`,
            },
        });
        
        if (!response.ok) {
            throw new Error('Failed to close position');
        }
    }
}

const apiService = new ApiService();

export default apiService;
