import { login, refreshToken as apiRefreshToken } from "@/api";
import { User } from "@/types";
import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthStore {
    user: User | null
    token: string | null;
    refreshToken: string | null;
    joinGame: (username: string) => Promise<void>;
    refreshAuthToken: () => Promise<string | null>;
    logout: () => void;
    isTokenExpired: (token: string) => boolean;
}

const useAuthStore = create<AuthStore>()(
    persist(
        (set, get) => ({
            user: null,
            token: null,
            refreshToken: null,
            joinGame: async (username: string) => {
                try {
                    const response = await login(username)
                    set({ user: response.user, token: response.token, refreshToken: response.refreshToken });
                } catch (error) {
                    console.error("Failed to join game:", error);
                    set({ user: null, token: null, refreshToken: null });
                }
            },
            refreshAuthToken: async (): Promise<string | null> => {
                const { refreshToken } = get();
                if (!refreshToken) {
                    console.error("No refresh token available");
                    set({ user: null, token: null, refreshToken: null });
                    return null;
                }

                try {
                    const response = await apiRefreshToken(refreshToken);
                    set({ 
                        user: response.user, 
                        token: response.token, 
                        refreshToken: response.refreshToken 
                    });
                    return response.token;
                } catch (error) {
                    console.error("Failed to refresh token:", error);
                    // Clear auth state if refresh fails
                    set({ user: null, token: null, refreshToken: null });
                    return null;
                }
            },
            logout: () => {
                set({ user: null, token: null, refreshToken: null });
            },
            isTokenExpired: (token: string): boolean => {
                if (!token) return true;
                
                try {
                    // Decode JWT payload (base64 decode the middle part)
                    const payload = JSON.parse(atob(token.split('.')[1]));
                    const currentTime = Date.now() / 1000; // Convert to seconds
                    
                    // Check if token is expired (with 30 second buffer)
                    return payload.exp < (currentTime + 30);
                } catch (error) {
                    console.error("Error checking token expiration:", error);
                    return true; // Consider invalid tokens as expired
                }
            }
        }),
        { name: "current-user-storage" }
    )
)

export default useAuthStore;
